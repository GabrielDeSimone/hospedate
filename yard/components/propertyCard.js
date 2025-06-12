import Link from "next/link";
import CurrencyBubble from "./currencyBubble";
import Spinner from "./spinner";

const PropertyCard = (props) => {

    const classes = [
        "bg-white",
        "rounded-xl",
        "shadow-lg ",
        "cursor-pointer",
        props.extraClasses || '',
    ]

    const preventIfLoading = (e, propStatus) => {
        if (propStatus === 'loading') {
            e.preventDefault()
        }
    }

    const getPropHref = () => {
        if (props.stayParams) {
            return `/properties/${props.id}?checkinDate=${props.stayParams.checkinDate}&checkoutDate=${props.stayParams.checkoutDate}&guests=${props.stayParams.guests}`
        } else {
            return `/properties/${props.id}`
        }
    }

    return (
        <div className={classes.join(' ')}>
            <Link href={getPropHref()} className={`flex flex-col ${props.status === 'loading' ? 'cursor-default': ''}`} onClick={(e) => {preventIfLoading(e, props.status)}}>
                {props.status === 'loading' ? (
                    <div className="h-[300px] flex flex-col justify-center items-center text-gray-400">
                        <p className="text-center">Propiedad en carga</p>
                        <p className="text-center">Id: {props.id}</p>
                        <Spinner extraClasses="mt-4 h-[50px] w-[50px]" />
                    </div>
                ) : (
                    <div>
                        <div className="flex justify-center items-center text-center h-[250px]"> {/* imagen */}
                            <div className="w-full h-[250px]">
                                <img src={props.images[0] + '?im_w=720'} className="object-cover h-full w-full rounded-t-xl" />
                            </div>
                        </div>
                        <div className="px-4 mt-3 truncate ...">{props.title}</div>
                        <div className="mt-3 flex flex-row h-[40px] mb-3">
                            <div className="h-full w-[50px] flex flex-row justify-center items-center ml-4"> {/* guests */}
                                <span className="material-symbols-outlined text-gray-600">group</span>
                                <span className="text-gray-500 ml-1">{props.maxGuests}</span>
                            </div>
                            <div className="flex flex-col justify-center h-full w-[38px] ml-auto mr-4 text-gray-500"> {/* currency */}
                                {/* <div className="ring-1 h-[17px] w-full rounded-md text-xs text-center ring-[#BDE9F2]">
                            ARS
                        </div> */}
                                <CurrencyBubble currency="USDT" />
                                {/*<div className="ring-1 h-[17px] w-full rounded-md text-xs text-center ring-[#BDE9F2]">*/}
                                {/*    USDT*/}
                                {/*</div>*/}
                            </div>
                            <div className="h-full w-[140px] flex flex-row-reverse mr-6"> {/* price */}
                                {props.stayParams && (
                                    <div className="h-full w-2/4 border-l-2 pl-3 ml-2"> {/* price total */}
                                        <div>
                                            ${props.price * props.stayParams.stayDays}
                                        </div>
                                        <div className="text-xs text-gray-500">
                                            {props.stayParams.stayDays} noches
                                        </div>
                                    </div>
                                )}
                                <div className="h-full w-2/4"> {/* price per night */}
                                    <div>
                                        ${props.price}
                                    </div>
                                    <div className="text-xs text-gray-500">
                                        por noche
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </Link>
        </div>
    )
}

export default PropertyCard