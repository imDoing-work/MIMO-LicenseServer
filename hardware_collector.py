#!/usr/bin/env python3
import os
import glob
import json
import subprocess
import re

# ============================================================
# 基础工具
# ============================================================

def read_file(path, default=""):
    try:
        with open(path) as f:
            return f.read().strip()
    except:
        return default


def normalize_pci(addr: str) -> str:
    addr = addr.lower()
    if addr.count(":") == 1:
        return "0000:" + addr
    return addr


# ============================================================
# Board UUID（优先 sysfs，fallback dmidecode）
# ============================================================

def get_board_uuid():
    uuid = read_file("/sys/class/dmi/id/product_uuid")
    if uuid:
        return uuid.lower()

    try:
        out = subprocess.check_output(
            "dmidecode -t system | grep 'UUID' | awk '{print $2}'",
            shell=True,
            stderr=subprocess.DEVNULL,
            text=True,
        )
        return out.strip().lower()
    except:
        return ""


# ============================================================
# MAC 地址（等价 getMacAddresses）
# ============================================================

def get_mac_addresses():
    macs = set()
    mac_re = re.compile(r"^([0-9a-f]{2}:){5}[0-9a-f]{2}$")

    for iface in os.listdir("/sys/class/net"):
        if iface == "lo":
            continue
        addr = read_file(f"/sys/class/net/{iface}/address").lower()
        if mac_re.match(addr):
            macs.add(addr)

    return sorted(macs)


# ============================================================
# Total Memory KB（等价 getTotalMemoryKB）
# ============================================================

def get_total_memory_kb():
    try:
        with open("/proc/meminfo") as f:
            for line in f:
                if line.startswith("MemTotal:"):
                    return int(line.split()[1])
    except:
        pass
    return 0


# ============================================================
# NVMe 扫描（kernel + SPDK 覆盖）
# ============================================================

def scan_nvme():
    devices = {}

    # ---------- Step 1: PCI NVMe ----------
    for pci in glob.glob("/sys/bus/pci/devices/*"):
        class_code = read_file(os.path.join(pci, "class"))
        if class_code != "0x010802":
            continue

        pci_addr = normalize_pci(os.path.basename(pci))
        vendor = read_file(os.path.join(pci, "vendor")).lower()

        devices[pci_addr] = {
            "vendor": vendor,
            "capacity": 0,
        }

        # ---------- kernel nvme ----------
        nvme_paths = glob.glob(os.path.join(pci, "nvme", "*"))
        if nvme_paths:
            nvme = nvme_paths[0]
            for ns in glob.glob(os.path.join(nvme, "nvme*n*")):
                sectors = int(read_file(os.path.join(ns, "size"), "0"))
                sector_size = int(read_file(os.path.join(ns, "queue/hw_sector_size"), "512"))
                if sectors > 0:
                    devices[pci_addr]["capacity"] = sectors * 512
                    break

    # ---------- Step 2: SPDK bdev（覆盖容量） ----------
    try:
        out = subprocess.check_output(
            ["/usr/local/mimo/SPDK_for_MIMO/scripts/rpc.py", "bdev_get_bdevs"],
            stderr=subprocess.DEVNULL,
            text=True,
        )
        bdevs = json.loads(out)
    except:
        bdevs = []

    for b in bdevs:
        nvme = b.get("driver_specific", {}).get("nvme")
        if not nvme:
            continue

        nv = nvme[0]
        pci = normalize_pci(nv.get("pci_address", ""))

        if pci in devices:
            cap = b.get("block_size", 0) * b.get("num_blocks", 0)
            if cap > 0:
                devices[pci]["capacity"] = cap

    return devices


# ============================================================
# 汇总 HardwareBind（最终输出）
# ============================================================

def collect_hardware_bind():
    nvmes = scan_nvme()

    vendors = set()
    total_cap = 0

    for d in nvmes.values():
        if d["vendor"]:
            vendors.add(d["vendor"])
        total_cap += d["capacity"]

    return {
        "board_uuid": get_board_uuid(),
        "mac_addresses": get_mac_addresses(),
        "nvme_vendors": sorted(vendors),
        "total_nvme_cap": total_cap,
        "nvme_count": len(nvmes),
        "total_memory_kb": get_total_memory_kb(),
    }


# ============================================================
# main
# ============================================================

if __name__ == "__main__":
    hw = collect_hardware_bind()
    print(json.dumps(hw, indent=2))
