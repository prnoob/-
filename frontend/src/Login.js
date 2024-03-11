import React, { Component } from 'react';
import {BrowserRouter as Router, Switch, Route, Redirect, Routes, Navigate} from 'react-router-dom';
import ChatRoom from './ChatRoom';
import LoginForm from './LoginForm';
import RegisterForm from './RegisterForm';
import axios from 'axios';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            name: '',
            email: '',
            password: '',
            showLogin: true, // 默认显示登录表单
        };
    }

    // handleChange = (e) => {
    //     this.setState({ [e.target.name]: e.target.value });
    // };
    //
    // handleRegister = (e) => {
    //     e.preventDefault();
    //     axios
    //         .post('http://localhost:8080/register', {
    //             name: this.state.name,
    //             email: this.state.email,
    //             password: this.state.password,
    //         })
    //         .then((response) => {
    //             console.log(response.data);
    //         })
    //         .catch((error) => {
    //             console.error(error);
    //         });
    // };
    //
    // handleLogin = (e) => {
    //     e.preventDefault();
    //     axios
    //         .post('http://localhost:8080/login', {
    //             name: this.state.name,
    //             email: this.state.email,
    //             password: this.state.password,
    //         })
    //         .then((response) => {
    //             console.log(response.data);
    //
    //             // 存储会话标识符(例如 JWT)
    //             const sessionId = response.data.sessionId;
    //             localStorage.setItem('sessionId', sessionId);
    //
    //
    //             // 重定向到聊天室界面
    //             return <Navigate to="/chat" replace />; // Assuming redirection after login
    //         })
    //         .catch((error) => {
    //             console.error(error);
    //         });
    // };

    // toggleForm = () => {
    //     this.setState((prevState) => ({
    //         showLogin: !prevState.showLogin,
    //     }));
    // };

    render() {
        const { showLogin, name, email, password } = this.state;

        return (
            <Router>
                <Routes>
                    <Route path="/" element={<LoginForm />} />
                    <Route path="/login" element={<LoginForm />} />
                    <Route path="/register" element={<RegisterForm />} />
                    <Route path="/chat" element={<ChatRoom />} />
                </Routes>
            </Router>
        );
    }
}

export default App;