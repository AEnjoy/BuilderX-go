import axios from 'axios'

let Port='18088'
function validateStatus(status) {
	return status<=500
}

const API = axios.create({
	baseURL:'http://localhost:'+Port, //
	timeout: 2000,                   //ms
	withCredentials:true,
	validateStatus
})

export default API