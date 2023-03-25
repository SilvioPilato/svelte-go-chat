export type WSMessage = {
    topic: "error" | "chat" | "login_success" | "login_user" | "logout_user"
    message: string
    username?: string
}

export type LoggedUser = {
    fullUsername: string;
}