import { render, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import AppBar from './AppBar';

describe('AppBar', () => {
    beforeEach(() => {
        localStorage.setItem('auth', 'true');
        localStorage.setItem('taxi_role', 'user');
    });

    it('should render', () => {
        const { getByText } = render(
            <MemoryRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <AppBar />
            </MemoryRouter>
        );

        expect(getByText('My Orders')).toBeInTheDocument();
    })

    it('should render admin link when role is admin', () => {
        localStorage.setItem('taxi_role', 'admin');
        const { getByText } = render(
            <MemoryRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <AppBar />
            </MemoryRouter>
        );

        expect(getByText('Drivers')).toBeInTheDocument();
    });

    it('should redirect on logout', () => {
        const { getByText } = render(
            <MemoryRouter future={{ v7_startTransition: true, v7_relativeSplatPath: true }}>
                <AppBar />
            </MemoryRouter>
        );

        const logoutButton = getByText('Logout');
        expect(logoutButton).toBeInTheDocument();

        logoutButton.click();
    
        expect(localStorage.getItem('auth')).toBe('false');
        expect(localStorage.getItem('taxi_role')).toBeNull();

        waitFor(() => {
            expect(window.location.href).toContain('/login');
        })
    }); 
})