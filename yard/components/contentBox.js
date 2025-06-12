const ContentBox = (props) => {
    return (
      <div className="bg-white rounded-lg shadow-lg p-5 pb-10 w-[80%] sm:w-[600px] mt-[30px] sm:mt-[67px] ml-auto mr-auto">
            {props.children}
      </div>
    )
}

export default ContentBox