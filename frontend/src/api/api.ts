import { create } from "apisauce";
import { CONFIG } from "./../config";
import { AuthRequest, AuthResponse, UserChangeTimezoneRequest, VtuberSearchRequest, VtuberSelectRequest } from "../types/api";

const sauce = create({
  baseURL: CONFIG.BASE_URL,
  headers: { Accept: 'application/json' },
});

// Is not authorized by default. Should be updated with updateAuthorizationHeader before every usage
const sauceAuthorized = create({
  baseURL: CONFIG.BASE_URL,
  headers: {
    Accept: 'application/json',
  },
});

// Updates sauceAuthorized instance (puts new token into headers)
const updateAuthorizationHeader = (token: string) => {
  sauceAuthorized.setHeader('vtubergo-token', token);
};

const api = {
  authGet: function (params: AuthRequest) {
    return sauce.get<AuthResponse>(`/auth?${params}`);
  },
  vtuberSearchPost: function (params: VtuberSearchRequest) {
    return sauceAuthorized.post("/search", params);
  },
  vtuberOrgsGet: function () {
    return sauceAuthorized.get("/orgs");
  },
  vtuberSelectPost: function (params: VtuberSelectRequest) {
    return sauceAuthorized.post("/select", params);
  },
  userChangeTimezonePost: function (params: UserChangeTimezoneRequest) {
    return sauceAuthorized.post("/timezone", params);
  },
  userGetTimezoneGet: function () {
    return sauceAuthorized.get("/timezone");
  },
};

export { api, updateAuthorizationHeader };