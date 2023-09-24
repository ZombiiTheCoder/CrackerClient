import Config from "/configBuilder.js"

async function init() {

    await Embed_CopyFile("launcher/Launcher.jar", await RoamingAppdata()+"/.crackerClient/Launcher.jar")

    let vlist = document.getElementById("versions")
    const launcher = document.getElementById("launcher")
    const install = document.getElementById("InstallAndLaunch")
    const changeUser = document.getElementById("ChangeUser")
    const username = document.getElementById("username")
    let vNum = 0

    changeUser.addEventListener("click", async ()=>{ updateConfig() })
    launcher.addEventListener("change", async ()=>{ updateConfig() })
    vlist.addEventListener("change", async ()=>{ updateConfig() })
    install.addEventListener("click", async ()=>{
        // print(`javaw -jar ${await RoamingAppdata()}\\.crackerClient\\Launcher.jar`)
        await execute(`${await RoamingAppdata()}\\.crackerClient`, "javaw", `-jar`, `${await RoamingAppdata()}\\.crackerClient\\Launcher.jar`)
    })
    
    async function versionExist() {
        let vpath = "q"
        switch (launcher.options[launcher.selectedIndex].value) {
            case "Vanilla": vpath = `${await RoamingAppdata()}/.crackerClient/versions/${vlist.options[vlist.selectedIndex].value}`; break
            case "Fabric": vpath = `${await RoamingAppdata()}/.crackerClient/versions/Fabric_${vlist.options[vlist.selectedIndex].value}`; break
            case "Quilt": vpath = `${await RoamingAppdata()}/.crackerClient/versions/Quilt_${vlist.options[vlist.selectedIndex].value}`; break
            case "Forge": vpath = `${await RoamingAppdata()}/.crackerClient/versions/Forge_${vlist.options[vlist.selectedIndex].value}`; break
        }
        if (await FileExist(vpath)) { install.innerText = "Launch"; return true }
        install.innerText = "Install"; return false
    }

    async function mov(json) {
        const versions = json["versions"]
        vNum = versions.length
        for (let index = 0; index < vNum; index++) {
            const ver = document.createElement("option");
            ver.text = versions[index].id
            ver.value = versions[index].id
            vlist.append(ver)
        }
        versionExist()
        if (await FileExist(await RoamingAppdata()+"/.crackerClient/config.json")) {
            const json = JSON.parse(await ReadFile(await RoamingAppdata()+"/.crackerClient/config.json"))
            // Versions Index
            let x = document.getElementById("versions").querySelectorAll(`option[value="${json["version"]}"]`);
            if (x.length === 1) {
                console.log(x[0].index);
                vlist.selectedIndex = x[0].index;
            }
            // Launcher Index
            x = document.getElementById("launcher").querySelectorAll(`option[value="${json["launcher"]}"]`);
            if (x.length === 1) {
                console.log(x[0].index);
                launcher.selectedIndex = x[0].index;
            }
            username.value=json["name"]
        } else {
            await WriteFile(await RoamingAppdata()+"/.crackerClient/config.json", new Config("EnterUserNameHere", "1.20.2", "Vanilla").toString())
            username.value="EnterUserNameHere"
        }
    }
    
    if (await FileExist(await RoamingAppdata()+"/.crackerClient/versionManifest_v2.json")) {
        mov(JSON.parse(await ReadFile(await RoamingAppdata()+"/.crackerClient/versionManifest_v2.json")))
    } else {
        fetch("https://piston-meta.mojang.com/mc/game/version_manifest_v2.json").then(resp => {resp.json().then(mov(json))})
    }

    async function updateConfig() {
        versionExist()
        await WriteFile(await RoamingAppdata()+"/.crackerClient/config.json", new Config(username.value, vlist.options[vlist.selectedIndex].value, launcher.options[launcher.selectedIndex].value).toString())
    }

    document.getElementById("logo").src = await Embed_ReadFileAsDataUrl(await edition() == "dev" ? "imgs/"+await edition()+"_logo.png" : "imgs/logo.png");
}
init()