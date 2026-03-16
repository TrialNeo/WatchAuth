import shutil
import subprocess
from pathlib import Path

backend_path = Path("./backend")
frontend_path = Path("./frontend")
dst = Path("./dist")
backend_build_exe = "watchauth-backend.exe"

print(f"building backend,backend_path={backend_path}")
res = subprocess.run(
    args=["go", "build", "-ldflags", "-s -w", "-o", backend_build_exe],
    cwd=backend_path / "cmd",
    capture_output=True,
    text=True)
if res.returncode == 0:
    print("backend compilate successfully")
    print("now moving files...")
    shutil.move(backend_path / "cmd" / backend_build_exe, dst / backend_build_exe)
    config = Path(backend_path / "configs")
    config.mkdir(exist_ok=True)
    shutil.copyfile(config / "menus.json", dst / "configs" / "menus.json")
    shutil.copyfile(config / "config.yaml", dst / "configs" / "config.yaml")
    print("backend building successfully!")
else:
    print("backend building failed!")

print()

print(f"building frontend,frontend_path={frontend_path}")
res = subprocess.run(
    args=["npm.cmd", "run", "build-only"],
    cwd=frontend_path,
    capture_output=True,
    text=True
)
if res.returncode == 0:
    dst_frontend = dst / "frontend"
    dst_frontend.mkdir(exist_ok=True)
    shutil.copy(frontend_path / "dist", dst_frontend)
else:
    print("building frontend failed!")