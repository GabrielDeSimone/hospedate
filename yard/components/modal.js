import React from 'react';

function Modal({ show, children, onClose, position }) {
    if (!show) {
        return null;
    }
    const positions = {
        "top": "sm:align-top",
        "center": "sm:align-middle"
    }

    const containerClasses = `
    inline-block
    align-bottom
    bg-white
    rounded-lg
    text-left
    overflow-hidden
    shadow-xl
    transform
    transition-all
    sm:my-8
    ${positions[position || "center"]}
    sm:max-w-2xl
    sm:w-full`

    return (
        <div
            className="
                fixed
                z-10
                inset-0
                overflow-y-auto"
            aria-labelledby="modal-title"
            role="dialog"
            aria-modal="true"
            >
            <div
                className="
                    flex
                    items-end
                    justify-center
                    min-h-screen
                    pt-4
                    px-4
                    pb-20
                    text-center
                    sm:block
                    sm:p-0">
                <div
                    className="
                        fixed
                        inset-0
                        bg-gray-500
                        bg-opacity-75
                        transition-opacity
                        "
                    aria-hidden="true"
                    onClick={onClose}></div>
                <span
                    className="
                        hidden
                        sm:inline-block
                        sm:align-middle
                        sm:h-screen"
                    aria-hidden="true">&#8203;</span>
                <div className={containerClasses}>
                    <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                        {children}
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Modal;
