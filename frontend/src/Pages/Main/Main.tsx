import { useInitData, useThemeParams } from '@vkruglikov/react-telegram-web-app';
import { ConfigProvider, theme } from 'antd';
import { DispatchWithoutAction, useEffect, useState } from 'react';
import useBetaVersion from '../../hooks/useBetaVersion';
import { VtuberList } from '../../components/VtuberList/VtuberList';
import { VtuberFilters } from '../../components/VtuberFilters';
import { useDispatch, useSelector } from 'react-redux';
import { userSelector } from '../../logic/user/UserSelectors';
import { UserActions } from '../../logic/user/UserSlice';
import { DotLoading, Tabs } from 'antd-mobile';
import { VtuberActions } from '../../logic/vtuber/VtuberSlice';
import { SelectionType } from '../../types/api';
import { PAGE_SIZE } from '../../helpers/consts';
import { vtuberSelector } from '../../logic/vtuber/VtuberSelectors';
import { useSearchValue } from '../../hooks/useSearchValue';
import { VtuberShow } from '../../components/VtuberShow';
import { UserSettings } from '../../components/UserSettings';

const Main: React.FC<{
    onChangeTransition: DispatchWithoutAction
}> = ({ onChangeTransition }) => {
    const [colorScheme, themeParams] = useThemeParams();
    const [isBetaVersion, handleRequestBeta] = useBetaVersion(false);
    const [activeBtn, setActiveBtn] = useState(true);
    const [showWelcomeBack, setShowWelcomeBack] = useState(true);

    const [initDataUnsafe, initData] = useInitData();
    const { token } = useSelector(userSelector);

    const [searchVtuberName, vtuberName, setVtuberName] = useSearchValue("", 3e2);
    const [offset, setOffset] = useState(0);
    const [companiesIds, setCompaniesIds] = useState<number[]>([]);
    const [wavesIds, setWavesIds] = useState<number[]>([]);
    const [selectionName, setSelectionName] = useState<SelectionType>("all");

    const { page_meta } = useSelector(vtuberSelector);

    const [fetchNext, setFetchNext] = useState(true);
    const [hasMore, setHasMore] = useState(true);
    const resetFilters = () => {
        setVtuberName("");
        setCompaniesIds([]);
        setWavesIds([]);
        setSelectionName("all");
    }
    const resetSearchState = () => {
        dispatch(VtuberActions.setVtubers([]));
        setOffset(0);
        setFetchNext(true);
        setHasMore(true);
    }

    const [showTags, setShowTags] = useState(true);
    const [showPictures, setShowPictures] = useState(true);

    useEffect(() => {
        if (page_meta) {
            if (page_meta.page_size_resp < page_meta.page_size_req) {
                setHasMore(false);
            }
            setFetchNext(true);
        }
    }, [page_meta]);


    const loadMore = async () => {
        if (fetchNext) {
            console.log(offset);
            dispatch(VtuberActions.vtuberSearch.request({
                name: searchVtuberName,
                orgs: companiesIds,
                waves: wavesIds,
                selected: selectionName,
                offset: offset,
                page_size: PAGE_SIZE
            }));
            setOffset(offset + PAGE_SIZE);
            setFetchNext(false);
        }
    };

    const dispatch = useDispatch();
    useEffect(() => {
        dispatch(UserActions.setInitData(initData));
    }, [initData]);
    useEffect(() => {
        if (!token) {
            dispatch(UserActions.auth.request(initData ?? ""));
        }
    }, []);
    useEffect(() => {
        if (token) {
            dispatch(VtuberActions.vtuberOrgs.request());
        }
    }, [token]);

    return (
        <div className='w-full h-full'>
            <ConfigProvider
                theme={
                    themeParams.text_color
                        ? {
                            algorithm:
                                colorScheme === 'dark'
                                    ? theme.darkAlgorithm
                                    : theme.defaultAlgorithm,
                            token: {
                                colorText: themeParams.text_color,
                                colorPrimary: themeParams.button_color,
                                colorBgBase: themeParams.bg_color,
                            },
                        }
                        : undefined
                }
            >
                <div className="contentWrapper">
                    {!token ?
                        <div style={{ color: '#00b578' }}>
                            <DotLoading color='currentColor' />
                            <span>Авторизация</span>
                        </div>
                        :
                        <div className='flex flex-col gap-4'>
                            <Tabs>
                                <Tabs.Tab title="Filters" key="filters">
                                    <VtuberFilters
                                        key="filters"
                                        searchVtuberName={searchVtuberName}
                                        vtuberName={vtuberName}
                                        companiesIds={companiesIds}
                                        wavesIds={wavesIds}
                                        selectionName={selectionName}
                                        page={offset}
                                        setVtuberName={setVtuberName}
                                        setCompaniesIds={setCompaniesIds}
                                        setWavesIds={setWavesIds}
                                        setSelectionName={setSelectionName}
                                        setPage={setOffset}
                                        resetSearchState={resetSearchState}
                                    />
                                </Tabs.Tab>
                                <Tabs.Tab title="Show" key="show">
                                    <VtuberShow
                                        key="show"
                                        showTags={showTags}
                                        setShowTags={setShowTags}
                                        showPictures={showPictures}
                                        setShowPictures={setShowPictures}
                                    />
                                </Tabs.Tab>
                                <Tabs.Tab title="Settings" key="settings">
                                    <UserSettings
                                        key="settings"
                                    />
                                </Tabs.Tab>
                            </Tabs>
                            <VtuberList
                                resetSearchState={resetSearchState}
                                offset={offset}
                                setOffset={setOffset}
                                selectionName={selectionName}
                                loadMore={token ? loadMore : async () => { }}
                                hasMore={hasMore}
                                showTags={showTags}
                                showPictures={showPictures}
                            />
                        </div>}
                </div>
            </ConfigProvider>
        </div>
    );
};

export { Main };