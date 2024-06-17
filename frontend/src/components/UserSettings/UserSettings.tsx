import { Avatar, Flex, Switch, Typography } from 'antd';
import { Vtuber } from '../../types/types';
import { Divider, ImageViewer, Loading, Modal, Space, Tag } from 'antd-mobile';
import { useState } from 'react';

const UserSettings: React.FC = () => {

    const [shareState, setShareState] = useState(false);

    return (
        <Flex gap="middle" vertical>
            <Divider
                style={{
                    color: '#1677ff',
                    borderColor: '#1677ff',
                    borderStyle: 'dashed',
                }}
            >
                Statistics
            </Divider>
            123
            <Divider
                style={{
                    color: '#1677ff',
                    borderColor: '#1677ff',
                    borderStyle: 'dashed',
                }}
            >
                Settings
            </Divider>
            <Space direction='vertical' wrap className='items-start'>
                <Flex gap="small">
                    <Switch checked={shareState} onClick={() => {
                        if(!shareState) {
                            Modal.show({
                                title: "Share my profile",
                                content: "If you turn it on - other users can see your telegram profile in vtuber cards and contact you. It was made for possibility of discovering new users to share interest in vtubers of your choice. Proceed?",
                                closeOnAction: true,
                                actions: [
                                    {
                                        key: "yes",
                                        text: "Yes",
                                        primary: true,
                                        onClick: () => {setShareState(true);}
                                    },
                                    {
                                        key: "no",
                                        text: "No",
                                    }
                                ]
    
                            });
                        } else {
                            setShareState(false);
                        }
                    }}/>
                    <Typography.Text type='secondary'> Share my profile in vtuber cards</Typography.Text>
                </Flex>
                <Switch />
            </Space>
        </Flex>
    );
}

export { UserSettings };