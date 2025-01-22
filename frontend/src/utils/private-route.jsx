import React from 'react';
import { Route, Navigate, Outlet } from 'react-router-dom';

const PrivateRoute = ({ component: Component, ...rest }) => {
  const isLoggedIn = () => {
    if(localStorage.getItem('accountData')) {
      const token = JSON.parse(localStorage.getItem('accountData')).token;
      return token !== null;
    }
  };

  return isLoggedIn() ? <Outlet /> : <Navigate to="/login" />;
};

export default PrivateRoute;
