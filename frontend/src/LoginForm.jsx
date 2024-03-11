import React from 'react';

const LoginForm = ({ name, email, password, handleChange, handleLogin, toggleForm }) => {
    return (
        <>
            <h1>Login</h1>
            <form onSubmit={handleLogin}>
                <input
                    type="name"
                    name="name"
                    placeholder="Name"
                    value={name}
                    onChange={handleChange}
                />
                <input
                    type="email"
                    name="email"
                    placeholder="Email"
                    value={email}
                    onChange={handleChange}
                />
                <input
                    type="password"
                    name="password"
                    placeholder="Password"
                    value={password}
                    onChange={handleChange}
                />
                <button type="submit">Login</button>
            </form>
            <button onClick={toggleForm}>还没有注册?</button>
        </>
    );
};

export default LoginForm;