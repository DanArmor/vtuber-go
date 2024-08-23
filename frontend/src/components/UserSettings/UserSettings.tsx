import { Avatar, Flex, Switch, Typography } from 'antd';
import { Vtuber } from '../../types/types';
import { Divider, ImageViewer, Loading, Modal, Space, Tag, Slider, Stepper, Button } from 'antd-mobile';
import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { UserActions } from '../../logic/user/UserSlice';
import { userSelector } from '../../logic/user/UserSelectors';
import { useInitData } from '@vkruglikov/react-telegram-web-app';

const UserSettings: React.FC = () => {

    const dispatch = useDispatch();
    const { timezone_shift } = useSelector(userSelector);
    const [initDataUnsafe, initData] = useInitData();

    useEffect(() => {
        dispatch(UserActions.timezone_get.request());
    }, []);

    return (
        <Flex gap="middle" vertical>
            <Divider
                style={{
                    color: '#1677ff',
                    borderColor: '#1677ff',
                    borderStyle: 'dashed',
                }}
            >
                Timezone shift from UTC (in hours)
            </Divider>

            <Stepper
                max={14}
                min={-14}
                value={timezone_shift}
                onChange={value => {
                    dispatch(UserActions.timezone_change.request(value));
                }}
            />
            <Divider
                style={{
                    color: '#1677ff',
                    borderColor: '#1677ff',
                    borderStyle: 'dashed',
                }}
            >
                Advanced (If app doesn't work)
            </Divider>
            <Button
                color='primary'
                fill='solid'
                onClick={() => {
                    dispatch(UserActions.auth.request(initData ?? ""))
                }}
            >
                Get new token
            </Button>
        </Flex>
    );
}

export { UserSettings };