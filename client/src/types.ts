export type WSMessage = {
    topic: "login" | "error" | "chat" | "login_success"
    message: string
    username?: string
}