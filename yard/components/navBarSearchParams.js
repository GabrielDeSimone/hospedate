import { useContext } from "react";
import { SearchParamsContext } from "./searchParamsProvider";

const NavBarSearchParams = (props) => {

    const [searchParams, setSearchParams] = useContext(SearchParamsContext);

    const items = [
        {
            icon: "location_on",
            text: "default",
            searchParam: "city",
            format: identity
        },
        {
            icon: "date_range",
            text: "checkin",
            searchParam: "checkinDate",
            format: humanizeDate
        },
        {
            icon: "date_range",
            text: "checkout",
            searchParam: "checkoutDate",
            format: humanizeDate
        },
        {
            icon: "group",
            text: "guests",
            searchParam: "guests",
            format: guestsFormat
        }
    ]

    function getSearchParam(param) {
        if (searchParams) {
            return searchParams[param]
        } else {
            return param
        }
    }

    return (
        <ul className="
            flex
            flex-nowrap
            justify-between
            bg-red
            cursor-pointer
            hover:shadow-md
            transition
            duration-200
            ease-in-out
            rounded-xl
            h-11
            align-center"
            onClick={props.onClick}
        >
            {
                items.map(item => (
                    <li
                        key={item.searchParam}
                        className="
                            flex
                            px-4
                            align-center
                        ">
                        <span className="
                            material-symbols-outlined
                            text-2xl
                            bg-transparent
                            text-gray-700
                            pt-1
                        ">{item.icon}</span>
                        <div className="pt-2 pl-2 text-gray-500">
                            {item.format(getSearchParam(item.searchParam))}
                        </div>
                    </li>
                ))
            }
        </ul>
    )
}

function identity(x) {
    return x
}

function humanizeDate(dateStr) {
    const [year, month, day] = dateStr.split('-').map(Number);
    const date = new Date(Date.UTC(year, month - 1, day));
    const options = { day: 'numeric', month: 'long', timeZone: 'UTC' };
    return date.toLocaleDateString('es-ES', options);
}

function guestsFormat(guests) {
    return `${guests} personas`
}

export default NavBarSearchParams
