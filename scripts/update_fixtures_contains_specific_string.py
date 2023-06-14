import sys
import subprocess
from pathlib import Path

repo_dir = Path(__file__).parent.parent
directory = repo_dir / Path("test/integration/fixtures")

search_string = sys.argv[1]
failed_tests = []

for file in directory.glob("*"):
    if file.is_file():
        content = file.read_text()
        if search_string in content:
            print(f"Found {search_string} in {file}")
            test_case = file.name.split(".")[0]
            command = f"make ARGS=\"-run {test_case}\" fixtures"
            output = subprocess.run(command, shell=True, capture_output=True)
            if output.returncode == 0:
                print(f"Successfully ran '{command}'")
            else:
                print(f"Command {command} failed with error:")
                print(output.stderr.decode())
                failed_tests.append(test_case)
        else:
            print(f"Did not find {search_string} in {file}")

print(f"These fixture generations failed: {failed_tests}")
