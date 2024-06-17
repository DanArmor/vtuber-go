import {useState, useEffect, useRef} from 'react';

export type UseSearchValueHookReturnType = [
  string,
  string,
  (s: string) => void,
];

// Хук для создания значения (инпута), к которому будет подвязан запрос поиска
// Возвращает переменную состояния, которая обновляется не чаще чем раз в updateInterval -
// когда юзер перестал изменять значение инпута, а также пару переменная состояние - сеттер, для самого инпута
export const useSearchValue = (
  defaultValue: string,
  updateInterval: number = 8e2,
): UseSearchValueHookReturnType => {
  const [innerValue, setInnerValue] = useState(defaultValue);
  const [searchValue, setSearchValue] = useState(innerValue);

  const timerId = useRef<NodeJS.Timeout | number | null>(null);

  useEffect(() => {
    if (timerId?.current) clearTimeout(timerId.current as any);

    timerId.current = setTimeout(() => {
      setSearchValue(innerValue);
    }, updateInterval);
  }, [innerValue]);

  return [searchValue, innerValue, setInnerValue];
};
