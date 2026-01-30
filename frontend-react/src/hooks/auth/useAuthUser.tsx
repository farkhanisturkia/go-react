//import Js Cookie
import Cookies from 'js-cookie';

//interface User
interface User {
    id: number;
    name: string;
    username: string;
    email: string;
}

export const useAuthUser = (): User | null => {
    // Mengambil data user dari cookie
    const user = Cookies.get('user');

    // Jika ada data user, parse JSON dan kembalikan
    // Jika tidak ada, kembalikan null
    return user ? JSON.parse(user) as User : null;
};
