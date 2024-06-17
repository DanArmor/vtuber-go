export const setTokenToLocalStorage = (token: string) => {
    localStorage.setItem("vtubergo_token", token);
};

export const getTokenFromLocalStorage = () => {
    const token = localStorage.getItem("vtubergo_token");
    return token;
}

export const deleteTokenFromLocalStorage = () => {
    localStorage.removeItem("vtubergo_token");
}