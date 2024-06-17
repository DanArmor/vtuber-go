export type Vtuber = {
    id: number,
    youtube_channel_id: string,
    channel_name: string,
    english_name?: string,
    photo_url?: string,
    banner_url?: string,
    twitter?: string,
    video_count?: number,
    subscriber_count?: number,
    clip_count?: number,
    top_topics?: string[],
    inactive: boolean,
    twitch?: string,
    description: string,
    isSelected?: boolean,
    edges: {
        wave?: VtuberWave,
        users?: {id: number}[]
    }
}

export type VtuberWave = {
    id: number,
    name: string
    edges: {
        vtubers?: Vtuber[],
        org?: VtuberOrg
    }
}

export type VtuberOrg = {
    id: number,
    name: string,
    icon_url: string
    edges: {
        waves?: VtuberWave[]
    }
}

export type WebAppInitData = {
    query_id?: string
    user?: WebAppUser
    receiver?: WebAppUser
    chat?: WebAppChat
    start_param?: string
    can_send_after?: number
    auth_date: number
    hash: string
}

export type WebAppUser = {
    id: number
    is_bot?: boolean
    first_name: string
    last_name?: string
    username?: string
    language_code?: string
    is_premium?: true
    added_to_attachment_menu?: true
    allows_write_to_pm?: true
    photo_url?: string
}


export type WebAppChat = {
    id: number
    type: 'group' | 'supergroup' | 'channel'
    title: string
    username?: string
    photo_url?: string
}