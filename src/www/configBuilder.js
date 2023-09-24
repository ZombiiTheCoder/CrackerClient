export default class Config {
    name
    version
    launcher
    versionIndex
    launcherIndex
    versionAmount

    constructor(name, version, launcher){
        this.name = name
        this.version = version
        this.launcher = launcher
    }

    toString() {
        return `{
    "name":"${this.name}",
    "version": "${this.version}",
    "launcher": "${this.launcher}"
}`
    }

}