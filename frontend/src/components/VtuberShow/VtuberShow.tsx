import { useMemo, useState } from "react";
import { Flex } from 'antd';
import { Selector } from "antd-mobile";

const options = [
    {
        label: "Show tags",
        value: 1,
    },
    {
        label: "Show Pictures",
        value: 2,
    }
]

export type VtuberShowProps = {
    showTags: boolean
    setShowTags: (show: boolean) => void
    showPictures: boolean
    setShowPictures: (show: boolean) => void
}

const VtuberShow: React.FC<VtuberShowProps> = (props: VtuberShowProps) => {
    const {
        showTags,
        setShowTags,
        showPictures,
        setShowPictures
    } = { ...props };

    const [stateOptions, setStateOptions] = useState(options);
    useMemo(() => {
        setStateOptions(stateOptions.map((value) => {
            if (value.value === 1) {
                return { ...value, setter: setShowTags };
            } else if (value.value === 2) {
                return { ...value, setter: setShowPictures };
            } else {
                return value;
            }
        }));
    }, [setShowTags]);

    return (
        <Flex gap="middle" vertical>
            <Selector
                columns={2}
                options={options}
                defaultValue={[1, 2]}
                multiple
                onChange={(arr) => {
                    setShowTags(false);
                    setShowPictures(false);
                    arr.forEach((v) => {
                        switch (v) {
                            case 1: {
                                setShowTags(true);
                                break;
                            }
                            case 2: {
                                setShowPictures(true);
                                break;
                            }
                        }
                    });
                }}
            />
        </Flex>
    );
}

export { VtuberShow };