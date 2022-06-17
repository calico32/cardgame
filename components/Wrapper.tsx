import React, { HTMLProps } from 'react'

interface WrapperProps extends HTMLProps<HTMLDivElement> {
  page?: boolean
  childClassName?: string
}

const Wrapper: React.FC<WrapperProps> = ({
  children,
  className,
  page,
  childClassName,
  ...props
}) => {
  return (
    <div className={`max-w-[800px] mx-auto ${className ?? ''}`} {...props}>
      <div className={`${childClassName ?? ''} ${page ? 'mx-4 mt-5' : 'mx-4'}`}>{children}</div>
    </div>
  )
}

export default Wrapper
