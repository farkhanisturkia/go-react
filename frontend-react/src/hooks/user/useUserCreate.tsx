// import useMutation dari '@tanstack/react-query';
import { useMutation } from '@tanstack/react-query';

//import service API
import Api from '../../services/api';

// import js-cookie
import Cookies from 'js-cookie';

//interface UserRequest
interface UserRequest {
    name: string;
    username: string;
    email: string;
    password: string;
}

export const useUserCreate = () => {

    return useMutation({

        // mutation untuk create user
        mutationFn: async (data: UserRequest) => {

            //get token from cookies
            const token = Cookies.get('token');

            //menggunakan service API untuk register
            const response = await Api.post('/api/users', data, {
                headers: {
                    Authorization: `Bearer ${token}`,
                },
            });

            //mengembalikan response data
            return response.data;
        }
    });
};
