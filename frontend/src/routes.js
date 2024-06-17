import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import SignIn from './components/auth/signin';
import SignUp from './components/auth/signup';
import ResetPassword from './components/auth/resetpassword';
import Home from './pages/home';
import Voting from './pages/voting';
import Results from './pages/results';

const AppRoutes = () => (
    <Router>
        <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/signin" element={<SignIn />} />
            <Route path="/signup" element={<SignUp />} />
            <Route path="/resetpassword" element={<ResetPassword />} />
            <Route path="/voting" element={<Voting />} />
            <Route path="/results" element={<Results />} />
        </Routes>
    </Router>
);

export default AppRoutes;
