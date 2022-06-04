/**
 * @file 格式化日期
 */

import moment from 'moment';
import 'moment/dist/locale/zh-cn';

moment.locale('zh-cn');

export const formatTime = (time: number) => {
    return moment.unix(time).utc(true).format('YYYY/MM/DD HH:mm:ss');
};

export const formatRelativeTime = (time: number) => {
    var now = moment();
    if (((now / 1000) - time) < (24 * 60 * 60)) {
        return moment.unix(time).fromNow();
    }
    return moment.unix(time).format('YYYY/MM/DD HH:mm');
};

