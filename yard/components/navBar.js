import Link from 'next/link';
import Router from 'next/router'
import React, { useContext, useRef, useState, useEffect } from 'react';
import {UserContext, isGuest, isLoggedIn, isHost} from './userProvider';

import HospedateLogo from './hospedateLogo';
import NavBarSearchParams from './navBarSearchParams';
import Modal from './modal';
import SearchParamsWindow from './searchParamsWindow';
import LoadingBar from "./loadingBar";

function clickOutsideHelper(ref, buttonRef, handler) {
    useEffect(() => {
        function handleClickOutside(event) {
            if (ref.current && !ref.current.contains(event.target) &&
                !buttonRef.current.contains(event.target)) {
                handler()
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [ref]);
}

function handleDropdownOpener(event, dropdownOpen, setDropdownOpen) {
    event.preventDefault();
    setDropdownOpen(!dropdownOpen);
}

async function logout(event) {
    event.preventDefault();

    const logoutReq = await fetch('/api/logout')
    const result = await logoutReq.json()

    Router.reload();
}

const dropDownButtons = [
    { name: 'Registrarse', href: '/register', show: (isLoggedIn, isHost) => (!isLoggedIn) },
    { name: 'Iniciar sesión', href: '/login', show: (isLoggedIn, isHost) => (!isLoggedIn) },
    { name: 'Mis Reservas', href: '/reservations', show: (isLoggedIn, isHost) => (isLoggedIn) },
    { name: 'Anfitriones', href: '/hosts', show: (isLoggedIn, isHost) => (isLoggedIn && isHost) },
    { name: 'Ser anfitrión', href: '/hostsApplication', show: (isLoggedIn, isHost) => (isLoggedIn && !isHost) },
    { name: 'Precios', href: '/help/pricing', show: (isLoggedIn, isHost) => (isLoggedIn) },
    { name: 'Ayuda', href: '/help', show: (isLoggedIn, isHost) => (isLoggedIn) },
    { name: 'Cerrar sesión', href: '/api/logout', onclick: logout, show: (isLoggedIn, isHost) => (isLoggedIn) },
]


function NavBar(props) {
    const [dropdownOpen, setDropdownOpen] = useState(false);
    const [user, setUser] = useContext(UserContext);
    const dropDownRef = useRef(null);
    const dropdownOpenerRef = useRef(null);

    // searchParams variables
    let showModal, setShowModal;
    let handleModalClose, handleModalOpen;
    if (props.showSearchParams) {
        [showModal, setShowModal] = useState(false);

        handleModalClose = () => {setShowModal(false)}
        handleModalOpen = () => {setShowModal(true)}
    }


    clickOutsideHelper(dropDownRef, dropdownOpenerRef, () => {
        setDropdownOpen(false);
    });

    return (
        <nav className="bg-white border-b border-gray-200 px-2 sm:px-4 py-2.5 rounded relative">
            <div className="container flex flex-wrap items-center justify-between mx-auto h-16">
                <Link href="/" className="flex items-center">
                    {/* <img src="https://flowbite.com/docs/images/logo.svg" className="h-6 mr-3 sm:h-9" alt="Flowbite Logo" /> */}
                    {/* <span className="self-center text-xl font-semibold whitespace-nowrap dark:text-white">Hospedate</span> */}
                    <HospedateLogo />
                </Link>
                {/* Search bar */}
                {/* <div className="flex">
                    <button type="button" data-collapse-toggle="navbar-search" aria-controls="navbar-search" aria-expanded="false" className="md:hidden text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 rounded-lg text-sm p-2.5 mr-1" >
                        <svg className="w-5 h-5" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd"></path></svg>
                        <span className="sr-only">Search</span>
                    </button>
                    <div className="relative hidden md:block">
                        <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                            <svg className="w-5 h-5 text-gray-500" aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd"></path></svg>
                            <span className="sr-only">Search icon</span>
                        </div>
                        <input type="text" id="navbar-search" className="block w-full p-2 pl-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" placeholder="Search..." />
                    </div>
                </div> */}

                {props.showSearchParams && (<NavBarSearchParams onClick={handleModalOpen} />)}

                <div className="flex items-center">
                    {/* data-dropdown-toggle="user-dropdown"  data-dropdown-placement="bottom" */}
                    <button ref={dropdownOpenerRef} type="button" className="flex mr-3 text-sm bg-transparent rounded-full md:mr-0 focus:ring-4 focus:ring-gray-200" id="user-menu-button" aria-expanded="false" onClick={(event) => handleDropdownOpener(event, dropdownOpen, setDropdownOpen)}>
                        <span className="sr-only">Open user menu</span>
                        {/* <img className="w-8 h-8 rounded-full" src="/docs/images/people/profile-picture-3.jpg" alt="user photo" /> */}
                        {/* <img className="w-8 h-8 rounded-full" src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80" alt="user photo" /> */}
                        <span className="material-symbols-outlined w-8 h-8 text-3xl bg-transparent text-gray-700">account_circle</span>
                    </button>
                    {/* <!-- Dropdown menu --> */}
                    <div ref={dropDownRef} className={"absolute z-50 text-base list-none bg-white divide-y divide-gray-100 rounded-lg shadow w-[250px] mt-[70px] top-0 ml-[-200px] " + (dropdownOpen ? 'block' : 'hidden')} id="user-dropdown">
                        {
                            isLoggedIn(user) && (
                                <Link className="px-4 py-3 block hover:bg-gray-100" href="/me">
                                    <span className="block text-sm text-gray-900">{user.name}</span>
                                    <span className="block text-sm font-medium text-gray-500 truncate">Mi cuenta</span>
                                </Link>
                            )
                        }
                        
                        <ul className="py-2" aria-labelledby="user-menu-button">
                            {
                                dropDownButtons.map((dropDownButton) => (
                                    dropDownButton.show(isLoggedIn(user), isHost(user)) ? (
                                        <li key={dropDownButton.name}>
                                            <Link href={dropDownButton.href} onClick={dropDownButton.onclick} className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">{dropDownButton.name}</Link>
                                        </li>
                                    ) : null
                                ))
                            }
                        </ul>
                    </div>
                </div>
            </div>
            {props.showSearchParams && (
                <Modal show={showModal} onClose={handleModalClose} position="top">
                    <SearchParamsWindow onSubmit={handleModalClose} />
                </Modal>
            )}
            <LoadingBar />
        </nav>
    );
}

export default NavBar