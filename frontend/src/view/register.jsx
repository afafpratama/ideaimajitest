import React, { useState } from 'react';
import axios from 'axios';
import { Input, Button, Card, Typography } from '@material-tailwind/react';
import { useNavigate } from 'react-router-dom';

function Register() {
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [phone, setPhone] = useState('');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleRegister = async (event) => {
        event.preventDefault();
        try {
          const response = await axios.post('http://localhost:3000/register', {
            name,
            phone,
            username,
            password,
          });

          // Navigate to the login page
          navigate('/login');
        } catch (error) {
          console.error('Error logging in', error);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100" onSubmit={handleRegister}>
            <Card color="transparent" shadow={false}>
                <Typography variant="h4" color="blue-gray">
                    Register
                </Typography>
                <Typography color="gray" className="mt-1 font-normal">
                    Nice to meet you! Enter your details to register.
                </Typography>
                <form className="mt-8 mb-2 w-80 max-w-screen-lg sm:w-96">
                    <div className="mb-1 flex flex-col gap-6">
                        <Typography variant="h6" color="blue-gray" className="-mb-3">
                            Your Name
                        </Typography>
                        <Input
                            size="lg"
                            placeholder="John Doe"
                            className=" !border-t-blue-gray-200 focus:!border-t-gray-900"
                            labelProps={{
                                className: "before:content-none after:content-none",
                            }}
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                        />
                        <Typography variant="h6" color="blue-gray" className="-mb-3">
                            Your Username
                        </Typography>
                        <Input
                            size="lg"
                            placeholder="johndoe"
                            className=" !border-t-blue-gray-200 focus:!border-t-gray-900"
                            labelProps={{
                                className: "before:content-none after:content-none",
                            }}
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                        <Typography variant="h6" color="blue-gray" className="-mb-3">
                            Your Phone
                        </Typography>
                        <Input
                            size="lg"
                            placeholder="0899 9999 9999"
                            className=" !border-t-blue-gray-200 focus:!border-t-gray-900"
                            labelProps={{
                                className: "before:content-none after:content-none",
                            }}
                            value={phone}
                            onChange={(e) => setPhone(e.target.value)}
                        />
                        <Typography variant="h6" color="blue-gray" className="-mb-3">
                            Password
                        </Typography>
                        <Input
                            type="password"
                            size="lg"
                            placeholder="********"
                            className=" !border-t-blue-gray-200 focus:!border-t-gray-900"
                            labelProps={{
                                className: "before:content-none after:content-none",
                            }}
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <Button type="submit" className="mt-6" fullWidth>
                        sign up
                    </Button>
                    <Typography color="gray" className="mt-4 text-center font-normal" onClick={() => navigate('/login')} >
                        Already have an account?{" "}
                        <a href="#" className="font-medium text-gray-900">
                        Sign In
                        </a>
                    </Typography>
                </form>
            </Card>
        </div>
    )
  }
  
  export default Register