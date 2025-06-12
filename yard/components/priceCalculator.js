import utils from '../utils/utils'
import {TableH, TableHTd, TableHTh, TableHTr} from "./myTable";
import {useState} from "react";
import Link from "next/link";

const PriceCalculator = (props) => {

    const [nights, setNights] = useState(3)
    const subtotal = props.price * nights;
    const total_billed = Number(subtotal * (1 + utils.PROTECTED_RESERVATIONS_FEE_GUEST)).toFixed(2)
    const host_fee = Number(subtotal * utils.PROTECTED_RESERVATIONS_FEE_HOST).toFixed(2)
    const profit = Number(subtotal - host_fee).toFixed(2)
    const humanReadableGuestFee = Number(utils.PROTECTED_RESERVATIONS_FEE_GUEST * 100).toFixed(2);
    const humanReadableHostFee = Number(utils.PROTECTED_RESERVATIONS_FEE_HOST * 100).toFixed(2);

    return (
        <div className={`border mt-4 ${props.extraClasses}`}>
            <h3 className="mt-3 ml-3">Ganancia con una reserva protegida</h3>
            <TableH extraClasses="ml-0">
                <TableHTr>
                    <TableHTh>Cantidad de noches</TableHTh>
                    <TableHTd><input className="border" type="number" value={nights} onChange={(e) => {setNights(e.target.value)}} /></TableHTd>
                </TableHTr>
                <TableHTr>
                    <TableHTh>Subtotal por {nights} noches</TableHTh>
                    <TableHTd>$ {subtotal}</TableHTd>
                </TableHTr>
                <TableHTr>
                    <TableHTh>Precio a pagar por el huesped (+ {humanReadableGuestFee}%)</TableHTh>
                    <TableHTd>$ {total_billed}</TableHTd>
                </TableHTr>
                <TableHTr>
                    <TableHTh>Ganancia del anfitrión (- {humanReadableHostFee}%)</TableHTh>
                    <TableHTd>$ {profit}</TableHTd>
                </TableHTr>
            </TableH>
            <p className="text-sm text-gray-500 px-5 my-3">Nuestra comisión es del 50% menos que en otras plataformas (<Link target="_blank" className="text-hosp-light-blue hover:underline cursor-pointer" href="/help/faq">saber más</Link>).</p>
        </div>
    )
}

export default PriceCalculator