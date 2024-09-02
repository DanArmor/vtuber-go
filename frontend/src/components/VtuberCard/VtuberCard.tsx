import { Avatar, Flex, Typography } from 'antd';
import { Vtuber } from '../../types/types';
import { Divider, ImageViewer, Loading, Space, Tag } from 'antd-mobile';
import { useState } from 'react';

export type VtuberCardProps = {
    vtuber: Vtuber | undefined,
    setPopupVisible: (value: boolean) => void
}

const VtuberCard: React.FC<VtuberCardProps> = (props: VtuberCardProps) => {
    const {
        vtuber,
        setPopupVisible
    } = { ...props };


    const [showPhoto, setShowPhoto] = useState(false);
    const [showPhotoUrl, setShowPhotoUrl] = useState("");

    return (
        <Flex gap="middle" className='w-screen h-full p-10 tgBgColor' vertical rootClassName='overflow-y-auto'>
            {
                vtuber
                    ?
                    <Flex gap="unset" className='w-full h-full' vertical>
                        <ImageViewer
                            image={showPhotoUrl}
                            visible={showPhoto}
                            onClose={() => {
                                setShowPhoto(false);
                            }}
                        />
                        <Flex gap="small" className='w-full' vertical>
                            <Avatar
                                shape='circle'
                                className='self-center w-36 h-36 hover:cursor-pointer'
                                src={vtuber.photo_url}
                                onClick={() => {
                                    setShowPhotoUrl(vtuber.photo_url ?? "");
                                    setShowPhoto(true);
                                }}
                            />
                            <Flex className='items-center w-full' vertical>
                                <Typography.Title style={{ marginBottom: 0, marginTop: 0 }}>
                                    {vtuber.english_name}
                                </Typography.Title>
                                <Typography.Text type='secondary'>
                                    {vtuber.edges.wave?.edges.org?.name} / {vtuber.edges.wave?.name}
                                </Typography.Text>
                                <Space wrap>
                                    {vtuber.top_topics?.map((value) => (
                                        <Tag color='primary' fill='outline'>
                                            {value}
                                        </Tag>
                                    ))}
                                </Space>
                                {vtuber.inactive && <Tag className='mt-2' color='warning' fill='outline'>Graduated</Tag>}
                            </Flex>
                            <Typography.Text type='secondary'>
                                Youtube: <Typography.Link
                                    target='_blank'
                                    href={`https://www.youtube.com/channel/${vtuber.youtube_channel_id}`}
                                >
                                    {vtuber.channel_name}
                                </Typography.Link>
                            </Typography.Text>

                            <Typography.Text type='secondary'>
                                Twitter: <Typography.Link
                                    target='_blank'
                                    href={`https://twitter.com/${vtuber.twitter}`}
                                >
                                    {vtuber.twitter}
                                </Typography.Link>
                            </Typography.Text>

                            {vtuber.twitch && <Typography.Text type='secondary'>
                                Twitch: <Typography.Link
                                    target='_blank'
                                    href={`https://twitch.com/${vtuber.twitch}`}
                                >
                                    {vtuber.twitch}
                                </Typography.Link>
                            </Typography.Text>}

                            <Divider
                                style={{
                                    color: '#1677ff',
                                    borderColor: '#1677ff',
                                    borderStyle: 'dashed',
                                }}
                            >
                                Counters
                            </Divider>

                            <Typography.Text type='secondary'>
                                Video count: {vtuber.video_count}
                            </Typography.Text>
                            <Typography.Text type='secondary'>
                                Subscriber count: {vtuber.subscriber_count}
                            </Typography.Text>
                            <Typography.Text type='secondary'>
                                Clip count: {vtuber.clip_count}
                            </Typography.Text>
                        </Flex>
                    </Flex>

                    :
                    <Loading />
            }
        </Flex>
    );
}

export { VtuberCard };