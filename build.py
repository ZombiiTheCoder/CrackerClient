import os

def remove_file(file_path):
    if os.path.exists(file_path):
        os.remove(file_path)

def copy_file(source, destination):
    with open(source, "rb") as src_file, open(destination, "wb") as dest_file:
        dest_file.write(src_file.read())

# Compiling Dev - Go
print("Compiling Dev - Go")
os.chdir("src")

remove_file("www/favicon.ico")
remove_file("winres/icon.png")

copy_file("data/imgs/dev_logo.ico", "www/favicon.ico")
copy_file("data/imgs/dev_logo.png", "winres/icon.png")

os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres make")
os.system("go build -ldflags=\"-X main.edition=dev\" -o ../bin/Cracker-Client-Dev.exe")

# Compiling Stable - Go
print("Compiling Stable - Go")

remove_file("www/favicon.ico")
remove_file("winres/icon.png")

copy_file("data/imgs/logo.ico", "www/favicon.ico")
copy_file("data/imgs/logo.png", "winres/icon.png")

os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres make")
os.system("go build -ldflags=\"-H windowsgui\" -o ../bin/Cracker-Client.exe")

os.chdir("../bin")
os.system("Cracker-Client-Dev.exe")
