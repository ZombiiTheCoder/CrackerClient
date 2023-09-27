import { Config } from "/configBuilder.js"

const provider = document.getElementById("provider")
const signin = document.getElementById("signin_microsoft")
document.getElementById("home").addEventListener("click", async ()=>{location.href=`http://localhost:${await port()}/index.html`})

async function init() {
    const roamingAppdata = await RoamingAppdata();

    if (await FileExist(`${roamingAppdata}/.crackerClient/config.json`)) {
        const config = JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/config.json`));
        const authType = config["authType"]
        signin.style.display = authType == "microsoft" ? "block" : "none"
        provider.selectedIndex = provider.querySelectorAll(`option[value="${authType.toLowerCase()}"]`)[0].index
    }

    async function updateConfig() {
        signin.style.display = provider.options[provider.selectedIndex].value == "microsoft" ? "block" : "none"
        const z = JSON.parse(await ReadFile(`${roamingAppdata}/.crackerClient/config.json`))
        
        WriteFile(`${roamingAppdata}/.crackerClient/config.json`, new Config(z["name"], z["version"], z["launcher"], provider.options[provider.selectedIndex].value).toString());
    }

    provider.addEventListener("change", updateConfig);
    signin.addEventListener("click", ()=>{login_with_microsft()})
    document.getElementById("logo").src = await Embed_ReadFileAsDataUrl(await edition() == "dev" ? `imgs/${await edition()}_logo.png` : "imgs/logo.png");
}

init()
