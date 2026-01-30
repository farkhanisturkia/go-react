//import FC from react
import { FC, useState, useEffect, FormEvent } from "react";

//import SidebarMenu
import SidebarMenu from '../../../components/SidebarMenu';

//import useNavigate, useParams and Link from react-router
import { useNavigate, useParams, Link } from 'react-router';

//import custom hook useUserByById
import { useUserById } from '../../../hooks/user/useUserById';

//import custom hook useUserUpdate
import { useUserUpdate } from '../../../hooks/user/useUserUpdate';

//interface for validation errors
interface ValidationErrors {
    [key: string]: string;
}

const UserEdit: FC = () => {

    //initialize useNavigate
    const navigate = useNavigate();

    //initialize useParams
    const { id } = useParams();

    //define state user
    const [name, setName] = useState<string>('');
    const [username, setUsername] = useState<string>('');
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');

    //define state errors
    const [errors, setErrors] = useState<ValidationErrors>({});

    // inisialisasi useUSerById
    const { data: user } = useUserById(Number(id));

    //set data user to state
    useEffect(() => {
        
        if (user) {
            setName(user.name);
            setUsername(user.username);
            setEmail(user.email);
        }
    }, [user]);

    // Inisialisasi useUserUpdate
    const { mutate, isPending } = useUserUpdate();

    // Handle submit form
    const updateUser = async (e: FormEvent) => {
        e.preventDefault();

        // Call the user update mutation
        mutate({
            id: Number(id),
            data: {
                name,
                username,
                email,
                password
            }
        }, {
            onSuccess: () => {

                // Redirect to users index
                navigate('/admin/users');

            },
            onError: (error: any) => {
                
                //set errors to state "errors"
                setErrors(error.response.data.errors);

            }
        })
    }

    return (
        <div className="container mt-5 mb-5">
            <div className="row">
                <div className="col-md-3">
                    <SidebarMenu />
                </div>
                <div className="col-md-9">
                    <div className="card border-0 rounded-4 shadow-sm">
                        <div className="card-header">
                            EDIT USER
                        </div>
                        <div className="card-body">
                            <form onSubmit={updateUser}>

                                <div className="form-group mb-3">
                                    <label className="mb-1 fw-bold">Full Name</label>
                                    <input type="text" value={name} onChange={(e) => setName(e.target.value)} className="form-control" placeholder="Full Name" />
                                    {errors.Name && <div className="alert alert-danger mt-2 rounded-4">{errors.Name}</div>}
                                </div>

                                <div className="form-group mb-3">
                                    <label className="mb-1 fw-bold">Username</label>
                                    <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} className="form-control" placeholder="Username" />
                                    {errors.Username && <div className="alert alert-danger mt-2 rounded-4">{errors.Username}</div>}
                                </div>

                                <div className="form-group mb-3">
                                    <label className="mb-1 fw-bold">Email address</label>
                                    <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} className="form-control"
                                        placeholder="Email Address" />
                                    {errors.Email && <div className="alert alert-danger mt-2 rounded-4">{errors.Email}</div>}
                                </div>

                                <div className="form-group mb-3">
                                    <label className="mb-1 fw-bold">Password</label>
                                    <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="form-control"
                                        placeholder="Password" />
                                    {errors.Password && <div className="alert alert-danger mt-2 rounded-4">{errors.Password}</div>}
                                </div>

                                <button type="submit" className="btn btn-md btn-primary rounded-4 shadow-sm border-0" disabled={isPending}>
                                    {isPending ? 'Updating...' : 'Update'}
                                </button>

                                <Link to="/admin/users" className="btn btn-md btn-secondary rounded-4 shadow-sm border-0 ms-2">Cancel</Link>

                            </form>
                        </div>
                    </div>
                </div>
            </div>
        </div >
    )
}

export default UserEdit
