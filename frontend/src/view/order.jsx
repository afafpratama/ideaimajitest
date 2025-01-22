import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
  Card,
  CardHeader,
  Input,
  Typography,
  Button,
  CardBody,
  Chip,
  CardFooter,
  Dialog,
  DialogHeader,
  DialogBody,
  DialogFooter,
  Select,
  Textarea,
  IconButton,
} from "@material-tailwind/react";

// import Select from "react-select";
import { MagnifyingGlassIcon } from "@heroicons/react/24/outline";
import { XMarkIcon } from "@heroicons/react/24/outline";
 
const TABLE_HEAD = ["Customer", "Service", "Amount", ""];
 
const TABLE_ROWS = [];

function Order() {
    const [records, setRecords] = useState([]);
    const [filteredRecords, setFilteredRecords] = useState([]);
    const [search, setSearch] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const recordsPerPage = 10;

    const onLoad = async () => {
        await axios.get(`http://localhost:3000/order?${search ? `search=${search}&` : "" }page=${currentPage}&limit=${recordsPerPage}`, {
            headers: { 'X-JWT-TOKEN': JSON.parse(localStorage.getItem('accountData')).token }
        })
        .then(response => {
            setRecords(response.data.data);
            setFilteredRecords(response.data.data);
        })
        .catch(error => {
            console.error('Error fetching data:', error);
        });
    }

    // Fetch records from the API
    useEffect(() => {
        onLoad();
    }, [ ]);

    // Filter records based on the search term
    useEffect(() => {
        if (Array.isArray(records) && records.length > 0) {
            const filteredData = records.filter(record =>
                record.name.toLowerCase().includes(search.toLowerCase()) || record.service.toLowerCase().includes(search.toLowerCase())
            );
            setFilteredRecords(filteredData);
        } else {
            setFilteredRecords([]); // Ensure filteredRecords is always an array
        }
    }, [ records, search ]);

    // Pagination logic
    const indexOfLast = currentPage * recordsPerPage;
    const indexOfFirst = indexOfLast - recordsPerPage;
    const currentRecords = Array.isArray(filteredRecords) ? filteredRecords.slice(indexOfFirst, indexOfLast) : [];

    const paginate = (pageNumber) => setCurrentPage(pageNumber);

    const [customers, setCustomers] = useState([]);

    const [modalAdd, setModalAdd] = useState(false); 
    const handleModalAdd = () => setModalAdd(!modalAdd);

    const [newRecord, setNewRecord] = useState({ customer_id: '', service: '', amount: 0, unit: '', price: 0 });

    const handleSaveAdd = () => {
        newRecord.customer_id = Number(newRecord.customer_id);
        newRecord.amount = Number(newRecord.amount);
        newRecord.price = Number(newRecord.price);

        axios.post('http://localhost:3000/order', newRecord, {
            headers: { 'X-JWT-TOKEN': JSON.parse(localStorage.getItem('accountData')).token }
        })
        .then(response => {
            onLoad();
            setModalAdd(false);
        })
        .catch(error => {
            console.error('Error adding order:', error);
        });
      };

    const [selectedData, setSelectedData] = useState({ id: 0, customer_id: 0, name: '', phone: '', service: '', amount: 0, unit: '', price: 0 });

    const [modalEdit, setModalEdit] = useState(false); 
    const handleModalEdit = (record) => {
        setSelectedData(record);
        setModalEdit(!modalEdit);
    };

    const handleSaveEdit = () => {
        let editedData = { 
            id: Number(selectedData.id), 
            customer_id: Number(selectedData.customer_id), 
            service: selectedData.service, 
            amount: Number(selectedData.amount), 
            unit: selectedData.unit, 
            price: Number(selectedData.price) 
        }

        axios.put(`http://localhost:3000/order/${editedData.id}`, editedData, {
            headers: { 'X-JWT-TOKEN': JSON.parse(localStorage.getItem('accountData')).token }
        })
        .then(response => {
            onLoad();
            setModalEdit(false);
        })
        .catch(error => {
            console.error('Error updating order:', error);
        });
    };

    const [modalDelete, setModalDelete] = useState(false); 
    const handleModalDelete = (record) => {
        setSelectedData(record);
        setModalDelete(!modalDelete);
    };

    const handleConfirmDelete = () => {
        axios.delete(`http://localhost:3000/order/${selectedData.id}`, {
            headers: { 'X-JWT-TOKEN': JSON.parse(localStorage.getItem('accountData')).token }
        })
        .then(response => {
            setRecords(records.filter(record => record.id !== selectedData.id));
            setModalDelete(false);
        })
        .catch(error => {
            console.error('Error deleting order:', error);
        });
    };

    const [modalDetail, setModalDetail] = useState(false); 
    const handleModalDetail = (record) => {
        setSelectedData(record);
        setModalDetail(!modalDetail);
    };

    useEffect(() => {
        if( modalAdd || modalEdit ) {
            axios.get(`http://localhost:3000/customer?page=1&limit=100`, {
                headers: { 'X-JWT-TOKEN': JSON.parse(localStorage.getItem('accountData')).token }
            })
            .then(response => {
                setCustomers(response.data.data);
            })
            .catch(error => {
                console.error('Error fetching customer data:', error);
            });
        }
    }, [ modalAdd, modalEdit ]);

    const currencyFormatter = new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: 'IDR',
    });

    return (
        <div className="px-6 py-4">
            <div className="container mx-auto flex">
                <Card className="h-full w-full" shadow={false}>
                    <CardHeader floated={false} shadow={false} className="rounded-none">
                        <div className="mb-8 flex items-center justify-between gap-8">
                            <div>
                                <Typography variant="h5" color="blue-gray">
                                    Orders list
                                </Typography>
                                <Typography color="gray" className="mt-1 font-normal">
                                    See information about all orders
                                </Typography>
                            </div>
                            <div className="flex shrink-0 flex-col gap-2 sm:flex-row">
                                <Button className="flex items-center gap-3" variant="filled" size="md" onClick={handleModalAdd}>
                                    Add Order
                                </Button>
                            </div>
                        </div>
                        <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
                            <div className="w-full md:w-72">
                                <Input
                                    label="Search"
                                    icon={<MagnifyingGlassIcon className="h-5 w-5" />}
                                    value={search}
                                    onChange={(e) => setSearch(e.target.value)}
                                />
                            </div>
                        </div>
                    </CardHeader>
                    <CardBody className="overflow-scroll hide-scrollbar px-0">
                        <table className="mt-4 w-full min-w-max table-auto text-left">
                            <thead>
                                <tr>
                                    {TABLE_HEAD.map((head) => (
                                        <th
                                            key={head}
                                            className="border-y border-blue-gray-100 bg-blue-gray-50/50 p-4"
                                        >
                                        <Typography
                                            variant="small"
                                            color="blue-gray"
                                            className="font-normal leading-none opacity-70"
                                        >
                                            {head}
                                        </Typography>
                                        </th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody>
                                {currentRecords.map(
                                ({ id, customer_id, name, phone, service, amount, unit, price }, index) => {
                                    const isLast = index === TABLE_ROWS.length - 1;
                                    const classes = isLast
                                        ? "p-4"
                                        : "p-4 border-b border-blue-gray-50";
                    
                                    return (
                                        <tr key={id} className="hover:bg-gray-50">
                                            <td className={classes}>
                                                <div className="flex items-center gap-3">
                                                    <div className="flex flex-col">
                                                        <Typography
                                                            variant="small"
                                                            color="light-blue"
                                                            className="font-normal"
                                                            as="a" href="#" onClick={() => handleModalDetail({ id, customer_id, name, phone, service, amount, unit, price })}
                                                        >
                                                            {name}
                                                        </Typography>
                                                        <Typography
                                                            variant="small"
                                                            color="blue-gray"
                                                            className="font-normal opacity-70"
                                                        >
                                                            {phone}
                                                        </Typography>
                                                    </div>
                                                </div>
                                            </td>
                                            <td className={classes}>
                                                <div className="flex items-center gap-3">
                                                    <div className="flex flex-col">
                                                        <Typography
                                                            variant="small"
                                                            color="blue-gray"
                                                            className="font-normal"
                                                        >
                                                            {service}
                                                        </Typography>
                                                    </div>
                                                </div>
                                            </td>
                                            <td className={classes}>
                                                <div className="flex items-center gap-3">
                                                    <div className="flex flex-col">
                                                        <Typography
                                                            variant="small"
                                                            color="blue-gray"
                                                            className="font-normal"
                                                        >
                                                            {amount} {unit}
                                                        </Typography>
                                                    </div>
                                                </div>
                                            </td>
                                            <td className={classes}>
                                                <div className="flex items-center gap-3">
                                                    <Typography as="a" href="#" variant="small" color="blue-gray" className="font-medium" onClick={() => handleModalEdit({ id, customer_id, name, phone, service, amount, unit, price })}>
                                                        Update
                                                    </Typography>
                                                    <Typography as="a" href="#" variant="small" color="blue-gray" className="font-medium" onClick={() => handleModalDelete({ id, customer_id, name, phone, service, amount, unit, price })}>
                                                        Delete
                                                    </Typography>
                                                </div>
                                            </td>
                                        </tr>
                                    );
                                },
                                )}
                            </tbody>
                        </table>
                    </CardBody>
                    <CardFooter className="flex items-center justify-between border-t border-blue-gray-50 p-4">
                        <Typography variant="small" color="blue-gray" className="font-normal">
                            Page {currentPage} of {Math.ceil(filteredRecords.length / recordsPerPage)}
                        </Typography>
                        <div className="flex gap-2">
                            <Button variant="outlined" size="sm" disabled={currentPage === 1} onClick={() => paginate(currentPage - 1)}>
                                Previous
                            </Button>
                            <Button variant="outlined" size="sm" disabled={currentPage === Math.ceil(filteredRecords.length / recordsPerPage)} onClick={() => paginate(currentPage + 1)}>
                                Next
                            </Button>
                        </div>
                    </CardFooter>
                </Card>

                <Dialog open={modalAdd} handler={handleModalAdd} className="relative bg-white m-4 rounded-lg shadow-2xl text-blue-gray-500 antialiased font-sans text-base font-light leading-relaxed w-full md:w-2/3 lg:w-2/4 2xl:w-1/3 min-w-[80%] md:min-w-[66.666667%] lg:min-w-[50%] 2xl:min-w-[33.333333%] max-w-[80%] md:max-w-[66.666667%] lg:max-w-[50%] 2xl:max-w-[33.333333%] p-4">
                    <DialogHeader className="relative m-0 block">
                        <Typography variant="h4" color="blue-gray">
                            Add Order
                        </Typography>
                        <IconButton
                            size="sm"
                            variant="text"
                            className="!absolute right-3.5 top-3.5"
                            onClick={handleModalAdd}
                        >
                            <XMarkIcon className="h-4 w-4 stroke-2" />
                        </IconButton>
                    </DialogHeader>
                    <DialogBody className="space-y-4 pb-6" style={{ display: 'block' }}>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Customer
                            </Typography>
                            <select
                                className="peer w-full h-full font-sans font-normal text-left outline outline-0 focus:outline-0 disabled:bg-blue-gray-50 disabled:border-0 transition-all border text-sm px-3 py-2.5 rounded-[7px] border-blue-gray-200 !w-full !border-[1.5px] !border-blue-gray-200/90 !border-t-blue-gray-200/90 bg-white text-gray-800 ring-4 ring-transparent placeholder:text-gray-600 focus:!border-primary focus:!border-t-blue-gray-900 group-hover:!border-primary"
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={newRecord.customer_id}
                                onChange={(e) => {
                                    setNewRecord({ ...newRecord, customer_id: e.target.value })
                                }}
                            >
                                {customers.map((customer) => (
                                    <option key={customer.id} value={customer.id}>
                                        {customer.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Service
                            </Typography>
                            <Textarea
                                rows={7}
                                className="!w-full !border-[1.5px] !border-blue-gray-200/90 !border-t-blue-gray-200/90 bg-white text-gray-600 ring-4 ring-transparent focus:!border-primary focus:!border-t-blue-gray-900 group-hover:!border-primary"
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={newRecord.service} onChange={(e) => setNewRecord({ ...newRecord, service: e.target.value })}
                            />
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Amount
                            </Typography>
                            <Input
                                type="number"
                                inputMode="numeric"
                                color="gray"
                                size="lg"
                                name="name"
                                className="placeholder:opacity-100 focus:!border-t-gray-900"
                                containerProps={{
                                    className: "!min-w-full",
                                }}
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={newRecord.amount} onChange={(e) => setNewRecord({ ...newRecord, amount: e.target.value })}
                            />
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Unit
                            </Typography>
                            <Input
                                color="gray"
                                size="lg"
                                name="name"
                                className="placeholder:opacity-100 focus:!border-t-gray-900"
                                containerProps={{
                                    className: "!min-w-full",
                                }}
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={newRecord.unit} onChange={(e) => setNewRecord({ ...newRecord, unit: e.target.value })}
                            />
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Price (Rp)
                            </Typography>
                            <Input
                                type="number"
                                inputMode="numeric"
                                color="gray"
                                size="lg"
                                name="name"
                                className="placeholder:opacity-100 focus:!border-t-gray-900"
                                containerProps={{
                                    className: "!min-w-full",
                                }}
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={newRecord.price} onChange={(e) => setNewRecord({ ...newRecord, price: e.target.value })}
                            />
                        </div>
                    </DialogBody>
                    <DialogFooter>
                        <Button
                            variant="outlined"
                            color="blue-gray"
                            onClick={handleModalAdd}
                            className="mr-2"
                        >
                            <span>Cancel</span>
                        </Button>
                        <Button variant="gradient" color="light-blue" onClick={handleSaveAdd}>
                            <span>Add</span>
                        </Button>
                    </DialogFooter>
                </Dialog>

                <Dialog open={modalEdit} handler={handleModalEdit} className="relative bg-white m-4 rounded-lg shadow-2xl text-blue-gray-500 antialiased font-sans text-base font-light leading-relaxed w-full md:w-2/3 lg:w-2/4 2xl:w-1/3 min-w-[80%] md:min-w-[66.666667%] lg:min-w-[50%] 2xl:min-w-[33.333333%] max-w-[80%] md:max-w-[66.666667%] lg:max-w-[50%] 2xl:max-w-[33.333333%] p-4">
                    <DialogHeader className="relative m-0 block">
                        <Typography variant="h4" color="blue-gray">
                            Update Order
                        </Typography>
                        <IconButton
                            size="sm"
                            variant="text"
                            className="!absolute right-3.5 top-3.5"
                            onClick={handleModalEdit}
                        >
                            <XMarkIcon className="h-4 w-4 stroke-2" />
                        </IconButton>
                    </DialogHeader>
                    <DialogBody className="space-y-4 pb-6" style={{ display: 'block' }}>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Customer
                            </Typography>
                            <select
                                className="peer w-full h-full font-sans font-normal text-left outline outline-0 focus:outline-0 disabled:bg-blue-gray-50 disabled:border-0 transition-all border text-sm px-3 py-2.5 rounded-[7px] border-blue-gray-200 !w-full !border-[1.5px] !border-blue-gray-200/90 !border-t-blue-gray-200/90 bg-white text-gray-800 ring-4 ring-transparent placeholder:text-gray-600 focus:!border-primary focus:!border-t-blue-gray-900 group-hover:!border-primary"
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={selectedData.customer_id ? selectedData.customer_id.toString() : ""}
                                onChange={(e) => {
                                    setSelectedData({ ...selectedData, customer_id: e.target.value })
                                }}
                            >
                                {customers.map((customer) => (
                                    <option key={customer.id} value={customer.id}>
                                        {customer.name}
                                    </option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Service
                            </Typography>
                            <Textarea
                                rows={7}
                                className="!w-full !border-[1.5px] !border-blue-gray-200/90 !border-t-blue-gray-200/90 bg-white text-gray-600 ring-4 ring-transparent focus:!border-primary focus:!border-t-blue-gray-900 group-hover:!border-primary"
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={selectedData.service} onChange={(e) => setSelectedData({ ...selectedData, service: e.target.value })}
                            />
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Amount
                            </Typography>
                            <Input
                                type="number"
                                inputMode="numeric"
                                color="gray"
                                size="lg"
                                name="name"
                                className="placeholder:opacity-100 focus:!border-t-gray-900"
                                containerProps={{
                                    className: "!min-w-full",
                                }}
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={selectedData.amount} onChange={(e) => setSelectedData({ ...selectedData, amount: e.target.value })}
                            />
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Unit
                            </Typography>
                            <Input
                                color="gray"
                                size="lg"
                                name="name"
                                className="placeholder:opacity-100 focus:!border-t-gray-900"
                                containerProps={{
                                    className: "!min-w-full",
                                }}
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={selectedData.unit} onChange={(e) => setSelectedData({ ...selectedData, unit: e.target.value })}
                            />
                        </div>
                        <div>
                            <Typography
                                variant="small"
                                color="blue-gray"
                                className="mb-2 text-left font-medium"
                            >
                                Price (Rp)
                            </Typography>
                            <Input
                                type="number"
                                inputMode="numeric"
                                color="gray"
                                size="lg"
                                name="name"
                                className="placeholder:opacity-100 focus:!border-t-gray-900"
                                containerProps={{
                                    className: "!min-w-full",
                                }}
                                labelprops={{
                                    className: "hidden",
                                }}
                                value={selectedData.price} onChange={(e) => setSelectedData({ ...selectedData, price: e.target.value })}
                            />
                        </div>
                    </DialogBody>
                    <DialogFooter>
                        <Button
                            variant="outlined"
                            color="blue-gray"
                            onClick={handleModalEdit}
                            className="mr-2"
                        >
                            <span>Cancel</span>
                        </Button>
                        <Button variant="gradient" color="light-blue" onClick={handleSaveEdit}>
                            <span>Update</span>
                        </Button>
                    </DialogFooter>
                </Dialog>

                <Dialog open={modalDelete} handler={handleModalDelete} className="relative bg-white m-4 rounded-lg shadow-2xl text-blue-gray-500 antialiased font-sans text-base font-light leading-relaxed w-full md:w-2/3 lg:w-2/4 2xl:w-1/3 min-w-[80%] md:min-w-[66.666667%] lg:min-w-[50%] 2xl:min-w-[33.333333%] max-w-[80%] md:max-w-[66.666667%] lg:max-w-[50%] 2xl:max-w-[33.333333%] p-4">
                    <DialogHeader>Deletion of Order #{selectedData ? selectedData.id : ''}</DialogHeader>
                    <DialogBody style={{ display: 'grid' }}>
                        <div>
                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Name</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? selectedData.name : ''}</Typography>
                            </div>

                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Phone</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? selectedData.phone : ''}</Typography>
                            </div>

                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Service</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-small">{selectedData ? selectedData.service : ''}</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? `${selectedData.amount} ${selectedData.unit}` : ''}</Typography>
                            </div>

                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Price</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? currencyFormatter.format(selectedData.price) : ''}</Typography>
                            </div>
                        </div>
                        
                        <Typography variant="h6" color="blue-gray" className="my-2 text-left font-medium">Are you sure want to delete this record?</Typography>
                    </DialogBody>
                    <DialogFooter>
                        <Button
                            variant="outlined"
                            color="blue-gray"
                            onClick={handleModalDelete}
                            className="mr-2"
                        >
                            <span>Cancel</span>
                        </Button>
                        <Button variant="gradient" color="red" onClick={handleConfirmDelete}>
                            <span>Delete</span>
                        </Button>
                    </DialogFooter>
                </Dialog>

                <Dialog open={modalDetail} handler={handleModalDetail} className="relative bg-white m-4 rounded-lg shadow-2xl text-blue-gray-500 antialiased font-sans text-base font-light leading-relaxed w-full md:w-2/3 lg:w-2/4 2xl:w-1/3 min-w-[80%] md:min-w-[66.666667%] lg:min-w-[50%] 2xl:min-w-[33.333333%] max-w-[80%] md:max-w-[66.666667%] lg:max-w-[50%] 2xl:max-w-[33.333333%] p-4">
                    <DialogHeader>Detail of Order #{selectedData ? selectedData.id : ''}</DialogHeader>
                    <DialogBody style={{ display: 'grid' }}>
                        <div>
                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Name</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? selectedData.name : ''}</Typography>
                            </div>

                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Phone</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? selectedData.phone : ''}</Typography>
                            </div>

                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Service</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-small">{selectedData ? selectedData.service : ''}</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? `${selectedData.amount} ${selectedData.unit}` : ''}</Typography>
                            </div>

                            <div>
                                <Typography variant="small" color="blue-gray" className="mb-1 text-left font-medium">Price</Typography>
                                <Typography variant="small" color="blue-gray" className="mb-2 text-left font-small">{selectedData ? currencyFormatter.format(selectedData.price) : ''}</Typography>
                            </div>
                        </div>
                    </DialogBody>
                    <DialogFooter>
                        <Button variant="gradient" color="light-blue" onClick={handleModalDetail}>
                            <span>Close</span>
                        </Button>
                    </DialogFooter>
                </Dialog>
            </div>
        </div>
    )
  }
  
  export default Order