import NavBar from '../components/navBar';
import {UserContext, isGuest, isNotHost} from '../components/userProvider';
import { useContext, useEffect, useState } from 'react';
import Router from 'next/router'
import ContentCard from '../components/contentCard';
import PropertyGrid from '../components/propertyGrid';
import BalanceLabel from "../components/balanceLabel";
import Modal from "../components/modal";
import NewPropertyForm from "../components/newPropertyForm";
import toast from 'react-hot-toast';
import utils from "../components/utils";
import {TableV, TableVBody, TableVHeading, TableVTd, TableVThCol, TableVThRow, TableVTr} from "../components/myTable";
import Link from "next/link";
import SmartContainer from "../components/smartContainer";

async function fetchMyBalance() {
    const getBalance = await fetch(
        "/api/me/balance"
    )
    const response = await getBalance.json()

    if (response.data) {
        return response.data.amountCents
    } else {
        return null
    }
}

async function fetchMyProperties() {
    const getProperties = await fetch(
        "/api/me/properties"
    )
    const response = await getProperties.json()

    if (response.data) {
        return response.data
    } else {
        return []
    }
}

async function fetchMyPropsReservations() {
    const getMyPropsReservations = await fetch(
        "/api/me/properties/reservations"
    )
    const response = await getMyPropsReservations.json()
    if (response.data) {
        return response.data
    } else {
        toast.error("Hubo un error al cargar reservas de propiedades")
        return []
    }
}


const Hosts = (props) => {
    const [user, setUser] = useContext(UserContext);
    const [myProperties, setMyProperties] = useState([]);
    const [newPropertyModal, setNewPropertyModal] = useState(false);
    const [thereIsALoadingProperty, setThereIsALoadingProperty] = useState(false);
    const [loadingPropIntervalHandle, setLoadingPropIntervalHandle] = useState(null);
    const [myPropsReservations, setMyPropsReservations] = useState([])
    const [balanceCents, setBalanceCents] = useState(null);
    const handleNewPropertyModelOpen = () => {setNewPropertyModal(true)}
    const handleNewPropertyModelClose = () => {setNewPropertyModal(false)}
    const updateProps = async () => {
        const myProps = await fetchMyProperties();
        setThereIsALoadingProperty(myProps.map(prop => prop.status === 'loading').some(Boolean));
        setMyProperties(myProps);
    }
    const updateBalance = async () => {
        const myBalance = await fetchMyBalance()
        setBalanceCents(myBalance)
    }

    const closeModalAndReload = async () => {
        handleNewPropertyModelClose()
        await updateProps();
    }

    useEffect(() => {
        if (isGuest(user) || isNotHost(user)) {
            Router.push('/')
        }
    }, [user])

    useEffect(() => {
        // update properties
        updateProps();

        // update balance
        updateBalance();

        // update my props reservations
        (async () => {
            const reservations =  await fetchMyPropsReservations()
            setMyPropsReservations(reservations)
        })()
    }, [])

    useEffect(() => {
        if (thereIsALoadingProperty && !loadingPropIntervalHandle) {
            const handle = setInterval(() => {
                updateProps();
            }, 3000)
            setLoadingPropIntervalHandle(handle);
        } else if (!thereIsALoadingProperty && loadingPropIntervalHandle) {
            clearInterval(loadingPropIntervalHandle)
            setLoadingPropIntervalHandle(null);
            toast.success('Propiedad procesada', {
                id: 'propertyLoading'
            })
        }
    }, [thereIsALoadingProperty])


    return (
        <div>
            <NavBar />
            <SmartContainer yAxisSpaced={true} noTabletWidthLimit={true}>
                <div className="flex flex-col-reverse laptop:flex-row">
                    <div id="leftPanel" className="laptop:w-3/4 laptop:mr-4">
                        <div className="border-b border-gray-300 p-6">
                            <h2>Reservas de mis propiedades</h2>
                                {myPropsReservations.length === 0 && (
                                    <p className="text-gray-500">(No tienes reservas)</p>
                                )}
                                {myPropsReservations.length > 0 && (
                                    <TableV>
                                        <TableVHeading>
                                            <TableVThCol>Número de reserva</TableVThCol>
                                            <TableVThCol>Estado</TableVThCol>
                                            <TableVThCol>Fecha de check-in</TableVThCol>
                                            <TableVThCol>Fecha de check-out</TableVThCol>
                                            <TableVThCol>Precio total</TableVThCol>
                                            <TableVThCol>Ver más</TableVThCol>
                                        </TableVHeading>
                                        <TableVBody>
                                            {myPropsReservations.map((reservation) => (
                                                <TableVTr key={reservation.id}>
                                                    <TableVThRow># {reservation.id}</TableVThRow>
                                                    <TableVTd>{utils.humanizeReservationStatus(reservation.status)}</TableVTd>
                                                    <TableVTd>{utils.humanizeFullDate(reservation.checkinDate)}</TableVTd>
                                                    <TableVTd>{utils.humanizeFullDate(reservation.checkoutDate)}</TableVTd>
                                                    <TableVTd>$ {utils.humanizeCents(reservation.totalBilledCents)}</TableVTd>
                                                    <TableVTd><Link className="text-hosp-light-blue hover:underline cursor-pointer" href={`/reservations/${reservation.id}`}>Ver reserva</Link></TableVTd>
                                                </TableVTr>
                                            ))}
                                        </TableVBody>
                                    </TableV>
                                )}
                        </div>
                        <div className="relative p-6">
                            <h2>Mis propiedades</h2>
                            <a className="
                                text-hosp-light-blue
                                hover:underline
                                cursor-pointer
                                block
                                sm:absolute
                                sm:right-0
                                sm:top-0
                                sm:mr-6
                                sm:mt-6
                                "
                               onClick={handleNewPropertyModelOpen}
                            >Agregar nueva propiedad</a>
                            {myProperties.length > 0 && (
                                <PropertyGrid properties={myProperties} customGridClasses={[
                                    "grid",
                                    "grid-cols-1",
                                    "tablet:grid-cols-2",
                                    "laptop:grid-cols-2",
                                    "desktop:grid-cols-3",
                                    "gap-4",
                                ]} />
                            )}
                            {myProperties.length === 0 && (
                                <p className="text-gray-500">(No tienes propiedades)</p>
                            )}
                        </div>
                    </div>
                    <div id="rightPanel" className="laptop:w-1/4">
                        <ContentCard shadow={true} extraClasses="relative">
                            <h2>Saldo</h2>
                            <BalanceLabel currency="usdt" balance={utils.humanizeCents(balanceCents)} extraClasses="absolute right-0 top-0 mt-6 mr-4"/>
                        </ContentCard>
                    </div>
                </div>
            </SmartContainer>
            {/*<div className="flex flex-col md:flex-row w-full mt-4 px-4"></div>*/}
            <Modal show={newPropertyModal} onClose={handleNewPropertyModelClose} position="top">
                <NewPropertyForm onSubmit={closeModalAndReload} />
            </Modal>
        </div>
    )
}

export default Hosts
