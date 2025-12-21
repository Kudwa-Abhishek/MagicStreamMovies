//configuration code for axios to call out server side endpoints
import axios from 'axios';

const apiUrl = import.meta.env.VITE_API_BASE_URL;

//reading in the base components address for server side endpoints
export default axios.create({
    baseURL:apiUrl,
    headers:{'Content-Type':'application/json'}
})