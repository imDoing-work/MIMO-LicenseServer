#!/usr/bin/env python3
import os
import json
import subprocess
import re

# ============================================================
# 基础工具
# ============================================================

def run(cmd):
    try:
        return subprocess.check_output(
            cmd,
            stderr=subprocess.DEVNULL,
            text=True,
        )
    except Exception:
        return ""

def normalize_serial(s: str) -> str:
    return s.strip().upper()

# ============================================================
# Board UUID
# ============================================================

def get_board_uuid():
    try:
        with open("/sys/class/dmi/id/product_uuid") as f:
            return f.read().strip().lower()
    except:
        pass

    out = run(["sh", "-c", "dmidecode -t system | grep UUID | awk '{print $2}'"])
    return out.strip().lower()

# ============================================================
# MAC 地址
# ============================================================

def get_mac_addresses():
    macs = set()
    mac_re = re.compile(r"^([0-9a-f]{2}:){5}[0-9a-f]{2}$")

    for iface in os.listdir("/sys/class/net"):
        if iface == "lo":
            continue
        try:
            with open(f"/sys/class/net/{iface}/address") as f:
                mac = f.read().strip().lower()
                if mac_re.match(mac):
                    macs.add(mac)
        except:
            pass

    return sorted(macs)

# ============================================================
# Total Memory KB
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
# NVMe Serial — SPDK (Primary)
# ============================================================

def get_nvme_serials_spdk():
    serials = set()
    rpc = "/usr/local/mimo/SPDK_for_MIMO/scripts/rpc.py"

    out = run([rpc, "bdev_nvme_get_controllers"])
    if not out:
        return serials

    try:
        ctrlrs = json.loads(out)
    except:
        return serials

    for c in ctrlrs:
        name = c.get("name")
        if not name:
            continue

        info = run([
            rpc,
            "bdev_nvme_get_controller_health_info",
            "-c", name
        ])

        if not info:
            continue

        try:
            data = json.loads(info)
            serial = data.get("serial_number")
            if serial:
                serials.add(normalize_serial(serial))
        except:
            continue

    return serials

# ============================================================
# NVMe Serial — nvme-cli (Fallback)
# ============================================================

def get_nvme_serials_nvmecli():
    serials = set()

    out = run(["nvme", "list", "-o", "json"])
    if not out:
        return serials

    try:
        data = json.loads(out)
    except:
        return serials

    for dev in data.get("Devices", []):
        sn = dev.get("SerialNumber")
        if sn:
            serials.add(normalize_serial(sn))

    return serials

# ============================================================
# 汇总 HardwareBind
# ============================================================

def collect_hardware_bind():
    serials = get_nvme_serials_spdk()

    if not serials:
        serials = get_nvme_serials_nvmecli()

    return {
        "board_uuid": get_board_uuid(),
        "mac_addresses": get_mac_addresses(),
        "nvme_serials": sorted(serials),
        "total_memory_kb": get_total_memory_kb(),
    }

# ============================================================
# main
# ============================================================

if __name__ == "__main__":
    hw = collect_hardware_bind()
    print(json.dumps(hw, indent=2))
