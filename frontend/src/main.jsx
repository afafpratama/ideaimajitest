import { Switch, ThemeProvider } from '@material-tailwind/react'
import React from 'react'
import ReactDOM from 'react-dom/client'
import { Navigate, Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import './index.css'
import 'tailwindcss/tailwind.css'

import Login from './view/login';
import Register from './view/register';
import Navbar from './component/navbar';
import Order from './view/order';
import PrivateRoute from './utils/private-route';
import Customer from './view/customer';

ReactDOM.createRoot(document.getElementById('root')).render(
  <ThemeProvider>
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/login" />} />
        <Route path="/login" element={<Login/>}/>
        <Route path="/register" element={<Register/>}/>
        <Route element={<PrivateRoute />}>
          <Route path="/order" element={
            <div className="relative mb-8 h-full w-full bg-white">
              <Navbar />
              <Order />
            </div>
          } />
        </Route>
        <Route element={<PrivateRoute />}>
          <Route path="/customer" element={
            <div className="relative mb-8 h-full w-full bg-white">
              <Navbar />
              <Customer />
            </div>
          } />
        </Route>
        {/* <Route element={<PrivateRoute />}>
          <Route path="/account" element={
            <div className="relative mb-8 h-full w-full bg-white">
              <Navbar />
              <Account />
            </div>
          } />
        </Route> */}
      </Routes>
    </Router>
  </ThemeProvider>
)
