import { useState } from 'react';
import {
  ColumnFiltersState,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  SortingState,
  useReactTable,
  VisibilityState,
} from '@tanstack/react-table';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { ChevronDown } from 'lucide-react';
import { Table } from '@/components/ui/table';
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import TableHeaderTransfer from '@/components/admin/transfer/table/table-header';
import {
  TableBodyTransfer,
  TableFooterTransfer,
  transferColumns,
} from '@/components/admin/transfer/table';
import { transfers } from '@/helpers/transfer_data';
import { AddTransfer } from '@/components/admin/transfer/modal/CreateModal';

export default function TransferPage() {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);
  const [columnVisibility, setColumnVisibility] = useState<VisibilityState>({});
  const [rowSelection, setRowSelection] = useState({});

  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 });

  const table = useReactTable({
    data: transfers,
    columns: transferColumns,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    onColumnVisibilityChange: setColumnVisibility,
    onRowSelectionChange: setRowSelection,
    state: {
      sorting,
      columnFilters,
      columnVisibility,
      rowSelection,
      pagination,
    },
    onPaginationChange: setPagination,
  });

  return (
    <div className="flex h-full overflow-hidden">
      <main className="flex-1 p-6 overflow-auto pb-20">
        <div className="flex-1 flex flex-col min-h-0">
          <Card className="w-full shadow-lg rounded-md border">
            <CardHeader className="p-4">
              <div className="flex justify-between items-center">
                <h3 className="text-lg font-semibold">Table Tansfer</h3>
                <div className="space-x-2">
                  <AddTransfer />
                </div>
              </div>
            </CardHeader>
            <CardContent className="p-4">
              <div className="flex items-center py-4">
                <Input
                  placeholder="Filter emails..."
                  value={
                    (table.getColumn('email')?.getFilterValue() as string) ?? ''
                  }
                  onChange={(event) =>
                    table.getColumn('email')?.setFilterValue(event.target.value)
                  }
                  className="max-w-sm"
                />
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="outline" className="ml-auto">
                      Columns <ChevronDown className="ml-2 h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    {table
                      .getAllColumns()
                      .filter((column) => column.getCanHide())
                      .map((column) => (
                        <DropdownMenuCheckboxItem
                          key={column.id}
                          className="capitalize"
                          checked={column.getIsVisible()}
                          onCheckedChange={(value: any) =>
                            column.toggleVisibility(!!value)
                          }
                        >
                          {column.id}
                        </DropdownMenuCheckboxItem>
                      ))}
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
              <div className="rounded-md border h-[525px] overflow-y-scroll">
                <Table>
                  <TableHeaderTransfer table={table} />
                  <TableBodyTransfer table={table} />
                </Table>
              </div>
            </CardContent>
            <CardFooter className="px-4 py-4 border-t">
              <TableFooterTransfer
                table={table}
                pagination={pagination}
                setPagination={setPagination}
              />
            </CardFooter>
          </Card>
        </div>
      </main>
    </div>
  );
}
