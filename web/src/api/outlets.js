import axios from 'axios';

// export const setLight = async (brand, deviceId, lightId, action) => {
//     const response = await axios.put(`/api/v1/device/light/${brand}/${deviceId}/${lightId}/${action}`);
//     return response.data;
// }

export const setOutlet = async (brand, id, action) => {
    const response = await axios.post(`/api/v1/device/outlet/${brand}/${id}/${action}`);
    return response.data;
}

export const getOutlets = async () => {
    try {
        const response = await fetch('/api/v1/device/outlet/kasa/test/discoverByKasa');
        const data = await response.json();
        // Return just the first response object which contains the ids
        return data;
    } catch (error) {
        console.error('Error fetching outlets:', error);
        throw error;
    }
}

export const getOutletState = async (brand, id) => {
    const response = await axios.post(`/api/v1/device/outlet/${brand}/${id}/state`);
    return response.data;
}

export const getOutletSysInfo = async (brand, id) => {
    const response = await axios.post(`/api/v1/device/outlet/${brand}/${id}/sysinfo`);
    return response.data;
}