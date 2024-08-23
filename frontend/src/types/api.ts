import { Vtuber, VtuberOrg, WebAppInitData } from "./types"

export type DefaultVtuberGoResult<T> = {
    result: T
}

export type VtuberGoError = {
    error: {
        code: number,
        msg: string
    }
}

export type AuthRequest = string;

export type SelectionType = "yes" | "no" | "all";

export type VtuberSearchRequest = {
    name?: string,
    orgs?: number[],
    waves?: number[],
    selected?: SelectionType,
    offset: number,
    page_size: number
};

export type VtuberSelectRequest = {
    vtuber_id: number
};

export type UserChangeTimezoneRequest = {
    timezone_shift: number
}

export type PageMetaInfo = {
    offset: number,
    page_size_req: number,
    page_size_resp: number
}

export type UserChangeTimezoneResponse = DefaultVtuberGoResult<{
    timezone_shift: number
}> & VtuberGoError;

export type UserGetTimezoneResponse = DefaultVtuberGoResult<{
    timezone_shift: number
}> & VtuberGoError;

export type AuthResponse = DefaultVtuberGoResult<{
    token: string
}> & VtuberGoError;

export type VtuberSearchResponse = DefaultVtuberGoResult<{
    vtubers: Vtuber[],
    page_meta: PageMetaInfo
}> & VtuberGoError;

export type VtuberOrgsResponse = DefaultVtuberGoResult<{
    orgs: VtuberOrg[]
}> & VtuberGoError;

export type VtuberSelectResponse = DefaultVtuberGoResult<{
    vtuber_id: number
    selected: boolean
}> & VtuberGoError;