import os
print("Compiling Dev - Go")
os.chdir("src")
try: os.remove("www/favicon.ico")
except: True
try: os.remove("winres/icon.png")
except: True
f=open("data/imgs/dev_logo.ico", "rb"); f2=open("www/favicon.ico", "wb"); f2.write(f.read()); f2.close(); f.close()
f=open("data/imgs/dev_logo.png", "rb"); f2=open("winres/icon.png", "wb"); f2.write(f.read()); f2.close(); f.close()
os.system("go install github.com/tc-hib/go-winres@latest")
# os.system("go-winres simply --icon \"data/imgs/dev_logo.png\" --manifest gui")
os.system("go-winres make")
os.system("go build -ldflags=\"-X main.edition=dev\" -o ../bin/Cracker-Client-Dev.exe")

print("Compiling Stable - Go")
try: os.remove("www/favicon.ico")
except: True
try: os.remove("winres/icon.png")
except: True
f=open("data/imgs/logo.ico", "rb"); f2=open("www/favicon.ico", "wb"); f2.write(f.read()); f2.close(); f.close()
f=open("data/imgs/logo.png", "rb"); f2=open("winres/icon.png", "wb"); f2.write(f.read()); f2.close(); f.close()
os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres make")
# os.system("go-winres simply --icon \"data/imgs/dev_logo.png\" --manifest gui")
os.system("go build -ldflags=\"-H windowsgui\" -o ../bin/Cracker-Client.exe")

os.chdir("../bin")
os.system("Cracker-Client-Dev.exe")