import NavBar from '../components/navBar';
import {isGuest, UserContext} from '../components/userProvider';
import {useContext, useEffect, useState} from 'react';
import utils from "../components/utils";
import Link from "next/link";
import {TableV, TableVBody, TableVHeading, TableVTd, TableVThCol, TableVThRow, TableVTr} from "../components/myTable";
import loadingBProvider from "../components/loadingBarProvider";
import ContentCard from "../components/contentCard";
import Router from "next/router";
import SmartContainer from "../components/smartContainer";

async function fetchMyReservations() {
    const getReservations = await fetch(
        "/api/me/reservations"
    )
    const response = await getReservations.json()

    if (response.data) {
        return response.data
    } else {
        return []
    }
}

const Reservations = (props) => {
    const [user, setUser] = useContext(UserContext);
    const [reservations, setReservations] = useState([])
    const [loadingBarStatus, setLoadingBarStatus] = useContext(loadingBProvider.LoadingBarContext);

    if (isGuest(user)) {
        Router.push('/')
    }

    useEffect(() => {
        (async () => {
            setLoadingBarStatus(loadingBProvider.LOADING_BAR_LOADING)
            const reservs = await fetchMyReservations()
            setReservations(reservs)
            setLoadingBarStatus(loadingBProvider.LOADING_BAR_READY)
        })()
    }, [])


    return (
        <div>
            <NavBar />
            <SmartContainer yAxisSpaced={true}>
                <ContentCard>
                    <h1 className="text-2xl mb-3">Mis reservas</h1>
                    {reservations.length > 0 ? (
                        <div className="relative overflow-x-auto">
                            <TableV>
                                <TableVHeading>
                                    <TableVThCol>Número de reserva</TableVThCol>
                                    <TableVThCol>Estado</TableVThCol>
                                    <TableVThCol>Fecha de check-in</TableVThCol>
                                    <TableVThCol>Fecha de check-out</TableVThCol>
                                    <TableVThCol>Precio total</TableVThCol>
                                    <TableVThCol>Forma de pago</TableVThCol>
                                    <TableVThCol>Ver más</TableVThCol>
                                </TableVHeading>
                                <TableVBody>
                                    {reservations.map((reservation) => (
                                        <TableVTr key={reservation.id}>
                                            <TableVThRow># {reservation.id}</TableVThRow>
                                            <TableVTd>{utils.humanizeReservationStatus(reservation.status)}</TableVTd>
                                            <TableVTd>{utils.humanizeFullDate(reservation.checkinDate)}</TableVTd>
                                            <TableVTd> {utils.humanizeFullDate(reservation.checkoutDate)}</TableVTd>
                                            <TableVTd>USDT $ {utils.humanizeCents(reservation.totalBilledCents)}</TableVTd>
                                            <TableVTd>{utils.humanizeReservationType(reservation.reservationType)}</TableVTd>
                                            <TableVTd><Link className="text-hosp-light-blue hover:underline cursor-pointer" href={`/reservations/${reservation.id}`}>Ver reserva</Link></TableVTd>
                                        </TableVTr>
                                    ))}
                                </TableVBody>
                            </TableV>
                        </div>
                    ) : (
                        <p className="text-gray-500">(No tienes reservas)</p>
                    )}
                </ContentCard>
            </SmartContainer>
        </div>
    )
}

export default Reservations
