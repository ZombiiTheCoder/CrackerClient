import { Config } from "/configBuilder.js"

const vlist = document.getElementById("versions")
const launcher = document.getElementById("launcher")
const install = document.getElementById("InstallAndLaunch")
const changeUser = document.getElementById("ChangeUser")
const username = document.getElementById("username")
const username_auth = document.getElementById("username_auth")
document.getElementById("settings").addEventListener("click", async ()=>{location.href=`http://localhost:${await port()}/settings.html`})

async function init() {
    const roamingAppdata = await RoamingAppdata();
    let LaunchCommand = ["-jar", `${roamingAppdata}/.crackerClient/Launcher.jar`]
    const launcherOptions = launcher.options;
    const vlistOptions = vlist.options;

    await Embed_CopyFile("launcher/Launcher.jar", `${roamingAppdata}/.crackerClient/Launcher.jar`);

    let vNum = 0;

    async function versionExist() {
        let vpath = "q";
        switch (launcherOptions[launcher.selectedIndex].value) {
            case "Vanilla":
                vpath = `${roamingAppdata}/.crackerClient/versions/${vlistOptions[vlist.selectedIndex].value}/${vlistOptions[vlist.selectedIndex].value}.jar`;
                break;
            case "Fabric":
                vpath = `${roamingAppdata}/.crackerClient/versions/Fabric_${vlistOptions[vlist.selectedIndex].value}/Fabric_${vlistOptions[vlist.selectedIndex].value}.jar`;
                break;
            case "Quilt":
                vpath = `${roamingAppdata}/.crackerClient/versions/Quilt_${vlistOptions[vlist.selectedIndex].value}/Quilt_${vlistOptions[vlist.selectedIndex].value}.jar`;
                break;
            case "Forge":
                vpath = `${roamingAppdata}/.crackerClient/versions/Forge_${vlistOptions[vlist.selectedIndex].value}/Forge_${vlistOptions[vlist.selectedIndex].value}.jar`;
                break;
        }
        install.innerText = await FileExist(vpath) ? "Launch" : "Install"

        if (await FileExist(`${roamingAppdata}/.crackerClient/config.json`) && await FileExist(`${roamingAppdata}/.crackerClient/AuthConfig.json`)) {
            const z = JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/config.json`))
            const z2 = JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/AuthConfig.json`))
            const authType = z["authType"]
            username_auth.value = z2["name"]
            changeUser.style.display = authType == "microsoft" ? "none" : "block"
            username.style.display = authType == "microsoft" ? "none" : "block"
            username_auth.style.display = authType == "microsoft" ? "block" : "none"
            LaunchCommand = authType == "microsoft" ? ["-cp", `${roamingAppdata}/.crackerClient/Launcher.jar`, "org.zombii.main.MainAuth"] : ["-jar", `${roamingAppdata}/.crackerClient/Launcher.jar`]
        }
    }

    async function updateConfig() {
        versionExist();
        WriteFile(`${roamingAppdata}/.crackerClient/config.json`, new Config(username.value, vlistOptions[vlist.selectedIndex].value, launcherOptions[launcher.selectedIndex].value, JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/config.json`))["authType"]).toString());
    }

    async function mov(json) {
        const versions = json["versions"];
        vNum = versions.length;
        for (let index = 0; index < vNum; index++) {
            const ver = document.createElement("option");
            ver.text = versions[index].id;
            ver.value = versions[index].id;
            vlist.append(ver);
        }
        
        versionExist();

        if (await FileExist(`${roamingAppdata}/.crackerClient/config.json`)) {
            const config = JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/config.json`));
            const name = config["name"]; const SelectedVersion = config["version"]; const SelectedLauncher = config["launcher"]
            vlist.selectedIndex = vlist.querySelectorAll(`option[value="${SelectedVersion}"]`)[0].index
            launcher.selectedIndex = launcher.querySelectorAll(`option[value="${SelectedLauncher}"]`)[0].index
            username.value = name;
        } else {
            WriteFile(`${roamingAppdata}/.crackerClient/config.json`, new Config("EnterUserNameHere", "1.20.2", "Vanilla", "cracked").toString());
            username.value = "EnterUserNameHere";
        }
    }

    if (await FileExist(`${roamingAppdata}/.crackerClient/versionManifest_v2.json`)) {
        mov(JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/versionManifest_v2.json`)));
    } else {
        fetch("https://piston-meta.mojang.com/mc/game/version_manifest_v2.json")
            .then(resp => resp.json())
            .then(json => mov(json));
    }

    changeUser.addEventListener("click", updateConfig);
    launcher.addEventListener("change", updateConfig);
    vlist.addEventListener("change", updateConfig);
    install.addEventListener("click", async () => {
        await execute(`${roamingAppdata}/.crackerClient`, "java", ...LaunchCommand);
    });
    document.getElementById("logo").src = await Embed_ReadFileAsDataUrl(await edition() == "dev" ? `imgs/${await edition()}_logo.png` : "imgs/logo.png");
    setTimeout(versionExist(), 100)
}

init()
