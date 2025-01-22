import React, { useState } from 'react';
import axios from 'axios';
import { Input, Button, Card, Typography } from '@material-tailwind/react';
import { useNavigate } from 'react-router-dom';

function Login() {
    const navigate = useNavigate();
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleLogin = async (event) => {
        event.preventDefault();

        try {
          const response = await axios.post('http://localhost:3000/login', {
            username,
            password,
          });

          // Extract account data from the response
        //   const { id, name, username, token } = response.data;

          // Create an array with the account data
          const accountData = { id: response.data.id, name: response.data.name, username: response.data.username, token: response.data.token };

          // Save the account data to local storage
          localStorage.setItem('accountData', JSON.stringify(accountData));

          // Navigate to the order page
          navigate('/order');

        //   console.log('Login successful', response.data);
        } catch (error) {
          console.error('Error logging in', error);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <Card color="transparent" shadow={false}>
                <Typography variant="h4" color="blue-gray">
                    Sign In
                </Typography>
                <Typography color="gray" className="mt-1 font-normal">
                    Nice to meet you! Enter your details to log in.
                </Typography>
                <form className="mt-8 mb-2 w-80 max-w-screen-lg sm:w-96" onSubmit={handleLogin}>
                    <div className="mb-1 flex flex-col gap-6">
                        <Typography variant="h6" color="blue-gray" className="-mb-3">
                            Username
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
                        sign in
                    </Button>
                    <Typography color="gray" className="mt-4 text-center font-normal" onClick={() => navigate('/register')} >
                        Don't have an account?{" "}
                        <a href="#" className="font-medium text-gray-900">
                        Sign Up
                        </a>
                    </Typography>
                </form>
            </Card>
        </div>
    )
  }
  
  export default Login