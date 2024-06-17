import { createRef, useEffect, useMemo, useState } from "react";
import { Vtuber, VtuberOrg } from "../../types/types";
import { Avatar, Checkbox, Flex, Input, InputRef, Select, Space, Switch, Typography } from 'antd';
import type { SelectProps } from "antd";
import { useDispatch, useSelector } from "react-redux";
import { vtuberSelector } from "../../logic/vtuber/VtuberSelectors";
import { SelectionType } from "../../types/api";
import { LabeledValue } from "antd/es/select";
import { SearchBar } from "antd-mobile";

const isSelectedOptions: SelectProps['options'] = [
    {
        label: "All",
        value: "all",
        desc: "All"
    },
    {
        label: "Only selected",
        value: "yes",
        desc: "Only selected"
    },
    {
        label: "Only not selected",
        value: "no",
        desc: "Only not selected"
    },
]

type OptionType = {
    id: number,
    label: string,
    value: string,
    icon_url?: string,
    desc: string
}

export type VtuberFiltersProps = {
    searchVtuberName: string
    vtuberName: string
    setVtuberName: (name: string) => void
    companiesIds: number[]
    setCompaniesIds: (ids: number[]) => void
    wavesIds: number[]
    setWavesIds: (ids: number[]) => void
    selectionName: SelectionType
    setSelectionName: (name: SelectionType) => void
    page: number
    setPage: (page: number) => void
    resetSearchState: () => void
}


const VtuberFilters: React.FC<VtuberFiltersProps> = (props: VtuberFiltersProps) => {
    const {
        searchVtuberName,
        vtuberName, setVtuberName,
        companiesIds, setCompaniesIds,
        wavesIds, setWavesIds,
        selectionName, setSelectionName,
        page, setPage,
        resetSearchState
    } = { ...props };

    const { orgs } = useSelector(vtuberSelector);
    const [orgsOptions, setOrgsOptions] = useState<SelectProps['options']>();
    const [wavesOptions, setWavesOptions] = useState<SelectProps['options']>();

    useEffect(() => {
        setOrgsOptions(
            orgs.map((value) => {
                return {
                    id: value.id,
                    label: value.name,
                    value: value.name,
                    desc: value.name,
                    icon_url: value.icon_url
                }
            })
        );

    }, [orgs]);
    useEffect(() => {
        setWavesOptions(            
            orgs.filter((val) => companiesIds.length === 0 ||!!companiesIds.find((v) => v === val.id)).reduce(
                (acc: OptionType[], val: VtuberOrg) => acc.concat(
                    val.edges.waves?.map((value) => {
                        return {
                            id: value.id,
                            label: value.name,
                            value: value.name,
                            desc: value.name
                        }
                    }) ?? []
                ),
                []
            )
        )
    }, [orgs, companiesIds]);

    useEffect(() => {
        resetSearchState();
    }, [searchVtuberName, companiesIds, wavesIds, selectionName]);

    return (
        <Flex gap="middle" vertical>
            <SearchBar
                placeholder="Vtuber name. . ."
                value={vtuberName}
                onChange={(v) => {
                    setVtuberName(v);
                }}
            />
            <Flex vertical gap="small">
                <Flex gap="small">
                    <Select
                        mode="multiple"
                        className="w-full"
                        placeholder="Company"
                        optionLabelProp="label"
                        options={orgsOptions}
                        optionRender={(option) => (
                            <Space>
                                {option.data.icon_url && <Avatar src={option.data.icon_url} />}
                                {option.data.desc}
                            </Space>
                        )}
                        onChange={(value, options) => {
                            setCompaniesIds(options.map((val: OptionType) => val.id));
                        }}
                    />
                    <Select
                        mode="multiple"
                        className="w-full"
                        placeholder="Wave"
                        optionLabelProp="label"
                        options={wavesOptions}
                        optionRender={(option) => (
                            <Space>
                                {option.data.desc}
                            </Space>
                        )}
                        onChange={(value, options) => {
                            setWavesIds(options.map((val: OptionType) => val.id));
                        }}
                    />
                </Flex>
                <Flex gap="small">

                    <Select
                        className="w-full"
                        defaultValue={"all"}
                        options={isSelectedOptions}
                        optionLabelProp="label"
                        placeholder="Is selected filter"
                        onChange={(value: SelectionType, option) => {
                            setSelectionName(value);
                        }}
                    />
                </Flex>

            </Flex>
        </Flex>
    );
}

export { VtuberFilters };