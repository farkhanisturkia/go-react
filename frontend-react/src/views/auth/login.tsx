// import FC from react
import { FC, useState, useContext, FormEvent } from 'react';

//import hook useNavigate from react router
import { useNavigate } from "react-router";

//import custom  hook useLogin from hooks
import { useLogin } from "../../hooks/auth/useLogin";

//import js-cookie
import Cookies from 'js-cookie'

//import context
import { AuthContext } from '../../context/AuthContext';

//interface for validation errors
interface ValidationErrors {
    [key: string]: string;
}

export const Login: FC = () => {

    //initialize navigate
    const navigate = useNavigate();

    //initialize useLogin
    const { mutate, isPending } = useLogin();

    //destruct auth context "setIsAuthenticated"
    const { setIsAuthenticated } = useContext(AuthContext)!;

    //define state
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');

    //define state for errors
    const [errors, setErrors] = useState<ValidationErrors>({});

    // Handle submit form
    const handleLogin = async (e: FormEvent) => {
        e.preventDefault();

        // Call the login mutation
        mutate({
            username,
            password
        }, {
            onSuccess: (data: any) => {
                const userData = data.data;

                if (userData.role?.toLowerCase() === 'admin') {
                    setErrors({
                        Error: 'Username or Password is incorrect'
                    });
                    return;
                }

                Cookies.set('token', userData.token);
                Cookies.set('user', JSON.stringify({
                    id: userData.id,
                    name: userData.name,
                    username: userData.username,
                    email: userData.email,
                    role: userData.role
                }));

                setIsAuthenticated(true);
                navigate('/admin/dashboard');
            },
            onError: (error: any) => {
                if (error) {
                    // setErrors(error.response.data.errors);
                    // setErrors(error.response.data.message);
                    setErrors({ Error: 'Username or Password is incorrect' });
                }
            }
        })
    }

    return (
        <div className="row justify-content-center mt-5">
            <div className="col-md-4">
                <div className="card border-0 rounded-4 shadow-sm">
                    <div className="card-body">
                        <h4 className='fw-bold text-center'>LOGIN</h4>
                        <hr />
                        {errors.Error && (
                            <div className="alert alert-danger mt-2 rounded-4">
                                {errors.Error}
                            </div>
                        )}
                        <form onSubmit={handleLogin}>
                            <div className="form-group mb-3">
                                <label className="mb-1 fw-bold">Username</label>
                                <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} className="form-control" placeholder="Username" />
                            </div>

                            <div className="form-group mb-3">
                                <label className="mb-1 fw-bold">Password</label>
                                <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="form-control" placeholder="Password" />
                            </div>
                            <button type="submit" className="btn btn-primary w-100 rounded-4" disabled={isPending}>
                                {isPending ? 'Loading...' : 'LOGIN'}
                            </button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Login;
