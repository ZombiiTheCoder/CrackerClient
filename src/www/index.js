async function init() {
    print(await wwwDev())
    if (await isEmbed()) {
        document.getElementById("logo").src = await ReadEmbedFileAsDataUrl("imgs/logo.png");
        return
    }
    document.getElementById("logo").src = await ReadFileAsDataUrl(await wwwDev()+"/imgs/logo.png");
}
init()