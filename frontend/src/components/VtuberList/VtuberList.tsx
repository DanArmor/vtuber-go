import { Avatar, Flex, Switch, Typography } from 'antd';
import { DotLoading, InfiniteScroll, Popup, PullToRefresh, Space, Tag } from "antd-mobile";
import { useDispatch, useSelector } from "react-redux";
import { vtuberSelector } from "../../logic/vtuber/VtuberSelectors";
import { VtuberActions } from '../../logic/vtuber/VtuberSlice';
import { PullStatus } from 'antd-mobile/es/components/pull-to-refresh';
import { useState } from 'react';
import { SelectionType } from '../../types/api';
import { UserOutlined } from '@ant-design/icons';
import { Vtuber } from '../../types/types';
import { VtuberCard } from '../VtuberCard';

const InfiniteScrollContent = ({ hasMore }: { hasMore?: boolean }) => {
    return (
        <>
            {
                hasMore ?
                    (
                        <>
                            <span>Loading</span>
                            <DotLoading />
                        </>
                    )
                    :
                    (
                        <span>No more data</span>
                    )
            }
        </>
    )
}

export type VtuberListProps = {
    loadMore: () => Promise<void>
    resetSearchState: () => void
    hasMore: boolean
    offset: number
    setOffset: (num: number) => void
    selectionName: SelectionType
    showTags: boolean
    showPictures: boolean
}

const VtuberList: React.FC<VtuberListProps> = (props: VtuberListProps) => {
    const { loadMore,
        hasMore,
        resetSearchState,
        offset,
        setOffset,
        selectionName,
        showTags,
        showPictures
    } = { ...props };


    const dispatch = useDispatch();
    const { vtubers } = useSelector(vtuberSelector);

    const [vtuberCardVisible, setVtuberCardVisible] = useState(false);
    const [vtuberCard, setVtuberCard] = useState<Vtuber | undefined>(undefined);

    const statusRecord: Record<PullStatus, string> = {
        pulling: 'Pull, if want to update',
        canRelease: 'Release - Let\'s update!',
        refreshing: 'Updating...',
        complete: 'Done!',
    }

    return (
        <PullToRefresh
            onRefresh={async () => {
                resetSearchState();
            }}
            renderText={status => {
                return <div>{statusRecord[status]}</div>
            }}
        >
            <Flex gap="middle" vertical>
                <Popup
                    position='left'
                    visible={vtuberCardVisible}
                    showCloseButton
                    mask
                    closeOnSwipe
                    onClose={() => {
                        setVtuberCardVisible(false);
                    }}
                >
                    <VtuberCard
                        vtuber={vtuberCard}
                        setPopupVisible={setVtuberCardVisible}
                    />
                </Popup>
                <Flex gap="middle" vertical>
                    {
                        vtubers.map((vtuber) => {
                            return (
                                <div className="flex items-center justify-between w-full">
                                    <div
                                        className="flex items-center hover:cursor-pointer"
                                        onClick={() => {
                                            setVtuberCard(vtuber);
                                            setVtuberCardVisible(true);
                                        }}
                                    >
                                        {showPictures &&
                                            <div>
                                                {vtuber.photo_url ?
                                                    <Avatar
                                                        size={'large'}
                                                        src={vtuber.photo_url}
                                                    />
                                                    :
                                                    <Avatar
                                                        size={'large'}
                                                        icon={<UserOutlined />}
                                                    />
                                                }
                                            </div>}
                                        <div className="flex flex-col items-start">
                                            <Typography.Text >{vtuber.english_name}</Typography.Text>
                                            <Typography.Text type="secondary">{vtuber.edges.wave?.edges.org?.name} / {vtuber.edges.wave?.name}</Typography.Text>
                                            {showTags &&
                                                <Space wrap>
                                                    {vtuber.top_topics?.map((value) => (
                                                        <Tag color='primary' fill='outline'>
                                                            {value}
                                                        </Tag>
                                                    ))}
                                                </Space>
                                            }
                                        </div>
                                    </div>
                                    <Switch checked={vtuber.isSelected} onClick={() => {
                                        if (selectionName === "yes") {
                                            if (vtuber.isSelected) {
                                                setOffset(offset - 1);
                                            } else {
                                                setOffset(offset + 1);
                                            }
                                        } else if (selectionName === "no") {
                                            if (vtuber.isSelected) {
                                                setOffset(offset + 1);
                                            } else {
                                                setOffset(offset - 1);
                                            }
                                        }
                                        dispatch(VtuberActions.vtuberSelect.request({ vtuber_id: vtuber.id }));
                                    }} />
                                </div>
                            );
                        })
                    }
                </Flex>
                <InfiniteScroll loadMore={loadMore} hasMore={hasMore}>
                    <InfiniteScrollContent hasMore={hasMore} />
                </InfiniteScroll>
            </Flex>
        </PullToRefresh>
    );
}

export { VtuberList };