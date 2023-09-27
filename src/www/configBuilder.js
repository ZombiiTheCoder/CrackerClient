export class Config {
    name
    version
    launcher
    authType

    constructor(name, version, launcher, authType){
        this.name = name
        this.version = version
        this.launcher = launcher
        this.authType = authType
    }

    toString() {
        return `{
    "name":"${this.name}",
    "version": "${this.version}",
    "launcher": "${this.launcher}",
    "authType": "${this.authType}"
}`
    }

}

export class AuthConfig {
    name
    uuid
    xuid
    clientID
    accessToken
    userType

    constructor(name, version, launcher){
        this.name = name
        this.version = version
        this.launcher = launcher
    }

    toString() {
        return `{
    "name":"${this.name}",
    "uuid": "${this.uuid}",
    "xuid": "${this.xuid}",
    "clientID": "${this.clientID}",
    "accessToken": "${this.accessToken}",
    "userType": "${this.userType}"
}`
    }
}