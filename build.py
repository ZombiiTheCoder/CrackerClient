import os
print("Compiling Dev - Go")
os.chdir("src")

os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres simply --icon \"data/imgs/dev_logo.png\" --manifest gui")
try: os.remove("www/favicon.ico")
except: True
f=open("data/imgs/dev_logo.ico", "rb"); f2=open("www/favicon.ico", "wb"); f2.write(f.read()); f2.close(); f.close()
os.system("go build -ldflags=\"-X main.edition=dev\" -o ../bin/Cracker_Client_Dev.exe")

print("Compiling Stable - Go")
os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres simply --icon \"data/imgs/logo.png\" --manifest gui")
try: os.remove("www/favicon.ico")
except: True
f=open("data/imgs/logo.ico", "rb"); f2=open("www/favicon.ico", "wb"); f2.write(f.read()); f2.close(); f.close()
os.system("go build -ldflags=\"-H windowsgui\" -o ../bin/Cracker_Client.exe")

os.chdir("../bin")
os.system("Cracker_Client_Dev.exe")