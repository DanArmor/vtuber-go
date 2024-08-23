import { restActionCreatorHelper } from './../../helpers/restActionCreatorHelper';

const userRestAction = restActionCreatorHelper(`user`);

export const AuthRestActions = userRestAction("auth");
export const TimezoneChangeRestActions = userRestAction("timezone_change");
export const TimezoneGetRestActions = userRestAction("timezone_get");
