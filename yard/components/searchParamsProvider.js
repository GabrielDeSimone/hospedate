import React, { useState } from 'react';

const SearchParamsProvider = (props) => {
    const [searchParams, setSearchParams] = useState(null);

    return (
        <SearchParamsContext.Provider value={[ searchParams, setSearchParams ]}>
            {props.children}
        </SearchParamsContext.Provider>
    );
}

export const SearchParamsContext = React.createContext(null);
export default SearchParamsProvider